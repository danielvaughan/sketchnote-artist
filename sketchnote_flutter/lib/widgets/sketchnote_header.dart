import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../theme/app_theme.dart';

class SketchnoteHeader extends StatelessWidget {
  const SketchnoteHeader({super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.palette_rounded, 
              color: AppTheme.accent, 
              size: 32, // Reduced icon size
            ),
            const SizedBox(width: 12),
            Text(
              'Sketchnote Artist',
              style: GoogleFonts.outfit(
                fontSize: 32, // Reduced font size
                fontWeight: FontWeight.w800,
                color: AppTheme.primaryText,
                letterSpacing: -1.0,
                height: 1.0, // Tighter line height
              ),
            ),
          ],
        ),
        const SizedBox(height: 8), // Reduced vertical spacing
        Text(
          'Transforming YouTube into Visual Knowledge',
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            fontWeight: FontWeight.w500,
            letterSpacing: 0.2, // Slightly tighter letter spacing
            fontSize: 15, // Slightly smaller subtitle
          ),
        ),
      ],
    );
  }
}
