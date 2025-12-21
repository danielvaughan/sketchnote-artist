import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../theme/app_theme.dart';

enum StatusType {
  idle,
  loading,
  success,
  error,
}

class StatusIndicator extends StatelessWidget {
  final String message;
  final StatusType type;

  const StatusIndicator({
    super.key,
    required this.message,
    required this.type,
  });

  @override
  Widget build(BuildContext context) {
    Color textColor;
    Color backgroundColor;
    IconData? icon;

    switch (type) {
      case StatusType.idle:
        textColor = AppTheme.secondaryText;
        backgroundColor = Colors.transparent;
        icon = Icons.info_outline_rounded;
        break;
      case StatusType.loading:
        textColor = AppTheme.primaryText;
        backgroundColor = const Color(0xFFE5E7EB);
        icon = null; // We'll show a spinner externally or standard icon
        break;
      case StatusType.success:
        textColor = const Color(0xFF1B5E20); // Green
        backgroundColor = const Color(0xFFE8F5E9);
        icon = Icons.check_circle_outline_rounded;
        break;
      case StatusType.error:
        textColor = const Color(0xFFB71C1C); // Red
        backgroundColor = const Color(0xFFFFEBEE);
        icon = Icons.error_outline_rounded;
        break;
    }

    // Modern "Chip" style
    return AnimatedContainer(
      duration: const Duration(milliseconds: 300),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(20),
        border: type != StatusType.idle
            ? Border.all(color: textColor.withValues(alpha: 0.1))
            : null,
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (type == StatusType.loading)
             SizedBox(
              width: 14,
              height: 14,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: textColor,
              ),
            )
          else if (icon != null)
            Icon(icon, size: 16, color: textColor),
            
          if (type != StatusType.idle) const SizedBox(width: 8),
          Text(
            message,
            style: GoogleFonts.inter(
              color: textColor,
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
