import 'dart:ui';
import 'package:flutter/material.dart';
import '../theme/app_theme.dart';
part 'sketchnote_input_button.dart';

class SketchnoteInput extends StatelessWidget {
  final TextEditingController controller;
  final VoidCallback onGenerate;
  final VoidCallback onClear;
  final bool isLoading;

  const SketchnoteInput({
    super.key,
    required this.controller,
    required this.onGenerate,
    required this.onClear,
    required this.isLoading,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      constraints: const BoxConstraints(maxWidth: 640),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(30),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
          child: Container(
            decoration: BoxDecoration(
              color: AppTheme.glassBase,
              borderRadius: BorderRadius.circular(30),
              border: Border.all(color: AppTheme.glassBorder, width: 1.5),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.08), // Increased shadow opacity
                  blurRadius: 24,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: controller,
                    onSubmitted: (_) => onGenerate(),
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppTheme.primaryText,
                          fontWeight: FontWeight.w500,
                        ),
                    decoration: InputDecoration(
                      hintText: 'Paste YouTube URL here...',
                      hintStyle: TextStyle(
                          color: AppTheme.secondaryText.withValues(alpha: 0.6)),
                      contentPadding: const EdgeInsets.symmetric(
                          horizontal: 28, vertical: 16),
                      border: InputBorder.none,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(right: 8.0),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      ValueListenableBuilder<TextEditingValue>(
                        valueListenable: controller,
                        builder: (context, value, child) {
                          final bool isNotEmpty = value.text.isNotEmpty;
                          return IconButton(
                            icon: Icon(
                              Icons.delete_sweep_rounded,
                              color: isNotEmpty
                                  ? AppTheme.secondaryText
                                  : AppTheme.secondaryText
                                      .withValues(alpha: 0.3),
                              size: 24,
                            ),
                            onPressed: isNotEmpty ? onClear : null,
                            tooltip: 'Clear input',
                          );
                        },
                      ),
                      const SizedBox(width: 4),
                      _AnimatedPlayButton(
                        isLoading: isLoading,
                        onPressed: onGenerate,
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
