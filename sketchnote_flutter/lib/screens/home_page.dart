import 'dart:async';
import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import 'package:path_provider/path_provider.dart';
import 'package:share_plus/share_plus.dart';
import 'package:flutter/foundation.dart';
import 'package:universal_html/html.dart' as html;
import 'package:youtube_player_flutter/youtube_player_flutter.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:google_fonts/google_fonts.dart';
import '../theme/app_theme.dart';
import '../services/api_service.dart';
import '../services/sse_service.dart';
import '../widgets/status_indicator.dart';
import '../widgets/sketchnote_header.dart';
import '../widgets/sketchnote_input.dart';
import '../widgets/dual_pane_layout.dart';
import '../widgets/sketchnote_result_card.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final TextEditingController _urlController = TextEditingController();
  final ApiService _apiService = ApiService(baseUrl: ''); // Relative to same host
  late final SseService _sseService = SseService(baseUrl: '');

  String _userId = 'local-user';
  String _statusMessage = 'Ready to create...';
  String? _imageUrl;
  String? _videoId;
  bool _isLoading = false;
  
  YoutubePlayerController? _youtubeController;
  StreamSubscription? _sseSubscription;

  StatusType _getStatusType() {
    if (_statusMessage.startsWith('Error')) return StatusType.error;
    if (_isLoading) return StatusType.loading;
    if (_imageUrl != null) return StatusType.success;
    return StatusType.idle;
  }

  @override
  void initState() {
    super.initState();
    _initUser();
  }

  Future<void> _initUser() async {
    final identity = await _apiService.getUserIdentity();
    setState(() {
      _userId = identity.email;
    });
  }

  void _generateSketchnote() async {
    final url = _urlController.text.trim();
    if (url.isEmpty) return;

    final videoId = YoutubePlayer.convertUrlToId(url);
    if (videoId == null) {
      setState(() {
        _statusMessage = 'Error: Please enter a valid YouTube URL';
        _isLoading = false;
      });
      return;
    }

    setState(() {
      _isLoading = true;
      _imageUrl = null;
      _videoId = videoId;
      _statusMessage = 'Initializing...';
      _youtubeController = YoutubePlayerController(
        initialVideoId: videoId,
        flags: const YoutubePlayerFlags(autoPlay: true, mute: false),
      );
    });

    try {
      final sessionId = await _apiService.createSession(_userId);
      
      _sseSubscription = _sseService.runSse(
        userId: _userId,
        sessionId: sessionId,
        message: url,
      ).listen((event) {
        if (event.error != null) {
          _handleError(event.error!);
          return;
        }

        _handleAgentEvent(event);
      }, onError: (e) {
        _handleError(e.toString());
      }, onDone: () {
        if (mounted) {
          setState(() {
            _isLoading = false;
          });
        }
      });

    } catch (e) {
      if (mounted) _handleError(e.toString());
    }
  }

  void _handleAgentEvent(dynamic event) {
    if (!mounted) return;
    if (event.models != null && event.models!.isNotEmpty) {
      setState(() => _statusMessage = 'Curator is analyzing video...');
    }

    final content = event.content;
    if (content != null && content['parts'] != null) {
      for (var part in content['parts']) {
        if (part['functionCall'] != null) {
          final toolName = part['functionCall']['name'] as String;
          if (toolName.contains('summarize')) {
            setState(() => _statusMessage = 'Summarizing video content...');
          } else if (toolName.contains('generate_image')) {
            setState(() => _statusMessage = 'Artist is sketching...');
          }
        }

        if (part['text'] != null) {
          final text = part['text'] as String;
          final regExp = RegExp(r'[\w\s-]+\.png');
          final match = regExp.firstMatch(text);
          if (match != null) {
            final filename = match.group(0)!.trim();
            setState(() {
              _imageUrl = '/images/$filename';
              _statusMessage = 'Sketchnote created successfully!';
              _isLoading = false;
            });
          }
        }
      }
    }
  }

  void _handleError(String message) {
    if (!mounted) return;
    setState(() {
      _statusMessage = 'Error: $message';
      _isLoading = false;
    });
  }

  void _resetUI() {
    setState(() {
      _urlController.clear();
      _imageUrl = null;
      _videoId = null;
      _isLoading = false;
      _statusMessage = 'Ready to create...';
      _youtubeController?.dispose();
      _youtubeController = null;
      _sseSubscription?.cancel();
      _sseSubscription = null;
    });
    _sseService.stopSse();
  }

  Future<void> _downloadImage() async {
    if (_imageUrl == null) return;
    try {
      if (kIsWeb) {
        // Web Download Logic
        final anchor = html.AnchorElement(href: _imageUrl!);
        anchor.download = _imageUrl!.split('/').last;
        anchor.style.display = 'none';
        html.document.body!.children.add(anchor);
        anchor.click();
        html.document.body!.children.remove(anchor);
        
        if (mounted) {
          setState(() => _statusMessage = 'Download started.');
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Image download started')),
          );
        }
        return;
      }

      // Mobile Download Logic
      setState(() => _statusMessage = 'Downloading image...');
      final dio = Dio();
      final tempDir = await getTemporaryDirectory();
      final filename = _imageUrl!.split('/').last;
      final savePath = '${tempDir.path}/$filename';

      await dio.download(_imageUrl!, savePath);

      if (mounted) {
        setState(() => _statusMessage = 'Image downloaded to temp directory.');
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Downloaded to: $savePath')),
        );
      }
    } catch (e) {
      if (mounted) _handleError('Download failed: $e');
    }
  }

  Future<void> _shareImage() async {
    if (_imageUrl == null) return;
    try {
      await Share.share('Check out this sketchnote: $_imageUrl');
    } catch (e) {
      // Share error
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          // Background Layer
          Positioned.fill(
            child: Container(
              decoration: const BoxDecoration(
                gradient: AppTheme.meshGradient,
              ),
            ),
          ),
          // Content Layer
          Positioned.fill(
            child: SingleChildScrollView(
              padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 20.0), // Reduced vertical padding
              child: Center(
                child: Container(
                  constraints: const BoxConstraints(maxWidth: 1100),
                  child: Column(
                    children: [
                      const SketchnoteHeader(),
                      const SizedBox(height: 20), // Reduced header-input gap
                      SketchnoteInput(
                        controller: _urlController,
                        onGenerate: _generateSketchnote,
                        onClear: _resetUI,
                        isLoading: _isLoading,
                      ),
                      const SizedBox(height: 32), // Reduced input-result gap
                      AnimatedSwitcher(
                        duration: const Duration(milliseconds: 600),
                        child: _videoId == null
                            ? _buildIdleState()
                            : _buildResultState(),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
          // Version Footer
          Positioned(
            bottom: 16,
            right: 24,
            child: FutureBuilder<PackageInfo>(
              future: PackageInfo.fromPlatform(),
              builder: (context, snapshot) {
                if (!snapshot.hasData) return const SizedBox.shrink();
                return Text(
                  'v${snapshot.data!.version}',
                  style: GoogleFonts.outfit(
                    color: AppTheme.secondaryText.withValues(alpha: 0.5),
                    fontSize: 12,
                    fontWeight: FontWeight.w500,
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildIdleState() {
    return Column(
      key: const ValueKey('idle'),
      children: [
        ShaderMask(
          shaderCallback: (bounds) => AppTheme.accentGradient.createShader(bounds),
          child: const Icon(
            Icons.auto_awesome_motion_rounded,
            size: 100,
            color: Colors.white, // Color must be white for ShaderMask to work
          ),
        ),
        const SizedBox(height: 24),
        Text(
          'Enter a YouTube URL to begin your visual journey',
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.bodyMedium,
        ),
      ],
    );
  }

  Widget _buildResultState() {
    return Column(
      key: const ValueKey('results'),
      children: [
        AnimatedOpacity(
          duration: const Duration(milliseconds: 400),
          opacity: 1.0,
          child: StatusIndicator(
            message: _statusMessage,
            type: _getStatusType(),
          ),
        ),
        const SizedBox(height: 24),
        DualPaneLayout(
          left: _youtubeController != null
              ? Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(24),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withValues(alpha: 0.1),
                        blurRadius: 30,
                        offset: const Offset(0, 10),
                      ),
                    ],
                  ),
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(24),
                    child: YoutubePlayer(controller: _youtubeController!),
                  ),
                )
              : AspectRatio(
                  aspectRatio: 16 / 9,
                  child: Container(
                    decoration: BoxDecoration(
                      color: Colors.black12,
                      borderRadius: BorderRadius.circular(24),
                    ),
                  ),
                ),
          right: SketchnoteResultCard(
            imageUrl: _imageUrl,
            onDownload: _downloadImage,
            onShare: _shareImage,
          ),
        ),
      ],
    );
  }

  @override
  void dispose() {
    _urlController.dispose();
    _youtubeController?.dispose();
    _sseSubscription?.cancel();
    _sseService.stopSse();
    super.dispose();
  }
}
