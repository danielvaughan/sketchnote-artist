import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../theme/app_theme.dart';

class SketchnoteResultCard extends StatelessWidget {
  final String? imageUrl;
  final VoidCallback onDownload;
  final VoidCallback onShare;

  const SketchnoteResultCard({
    super.key,
    this.imageUrl,
    required this.onDownload,
    required this.onShare,
  });

  @override
  Widget build(BuildContext context) {
    return AspectRatio(
      aspectRatio: 16 / 9,
      child: Container(
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
          child: Stack(
            children: [
              // Glass Background
              Positioned.fill(
                child: BackdropFilter(
                  filter: ImageFilter.blur(sigmaX: 5, sigmaY: 5),
                  child: Container(
                    decoration: BoxDecoration(
                      color: AppTheme.glassBase,
                      border: Border.all(color: AppTheme.glassBorder, width: 1.5),
                    ),
                  ),
                ),
              ),

              if (imageUrl != null)
                Positioned.fill(
                  child: AnimatedOpacity(
                    duration: const Duration(milliseconds: 500),
                    opacity: 1.0,
                    child: Image.network(
                      imageUrl!,
                      fit: BoxFit.contain,
                      loadingBuilder: (context, child, loadingProgress) {
                        if (loadingProgress == null) return child;
                        return const Center(child: CircularProgressIndicator());
                      },
                      errorBuilder: (context, error, stackTrace) {
                        return const Center(
                          child: Icon(Icons.broken_image, color: AppTheme.secondaryText, size: 48),
                        );
                      },
                    ),
                  ),
                ),
              if (imageUrl == null)
                const Center(
                  child: Icon(Icons.auto_awesome_rounded, color: AppTheme.secondaryText, size: 64),
                ),
              
              // Action Buttons - Always Visible (with subtle protection)
              Positioned(
                bottom: 20,
                right: 20,
                child: Container(
                  padding: const EdgeInsets.all(4), // Space for blur to bleed
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20),
                    boxShadow: [
                      BoxShadow(
                          color: Colors.black.withValues(alpha: 0.1),
                          blurRadius: 16,
                          spreadRadius: -4)
                    ],
                  ),
                  child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            _ModernActionButton(
                              icon: Icons.download_rounded,
                              label: 'Save',
                              onPressed: imageUrl != null ? onDownload : null,
                              isPrimary: true,
                            ),
                            const SizedBox(width: 12),
                            _ModernActionButton(
                              icon: Icons.share_rounded,
                              label: 'Share',
                              onPressed: imageUrl != null ? onShare : null,
                              isPrimary: false,
                            ),
                          ],
                        ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ModernActionButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback? onPressed;
  final bool isPrimary;

  const _ModernActionButton({
    required this.icon,
    required this.label,
    required this.onPressed,
    required this.isPrimary,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        boxShadow: isPrimary && onPressed != null
            ? [
                BoxShadow(
                  color: AppTheme.accent.withValues(alpha: 0.3),
                  blurRadius: 12,
                  offset: const Offset(0, 4),
                )
              ]
            : null,
      ),
      child: ElevatedButton.icon(
        onPressed: onPressed,
        icon: Icon(icon, size: 20),
        label: Text(
          label,
          style: const TextStyle(
            fontWeight: FontWeight.w600,
            letterSpacing: 0.3,
          ),
        ),
        style: ElevatedButton.styleFrom(
          backgroundColor: isPrimary ? AppTheme.accent : Colors.white,
          foregroundColor: isPrimary ? Colors.white : AppTheme.primaryText,
          disabledBackgroundColor: isPrimary 
              ? AppTheme.accent.withValues(alpha: 0.5) 
              : Colors.white.withValues(alpha: 0.5),
          disabledForegroundColor: Colors.white.withValues(alpha: 0.7),
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
          side: isPrimary ? null : const BorderSide(color: AppTheme.border),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
        ),
      ),
    );
  }
}
