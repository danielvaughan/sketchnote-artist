import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // Brand Colors
  static const Color background = Color(0xFFFBFBFE);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color primaryText = Color(0xFF1A1A1A);
  static const Color secondaryText = Color(0xFF6E6E73);
  static const Color accent = Color(0xFFFF3B30); // iOS-style red
  static const Color border = Color(0x1F000000);

  // Gradient Tokens
  static const LinearGradient meshGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [
      Color(0xFFFBFBFE),
      Color(0xFFF0F2F5),
      Color(0xFFE5E7EB),
    ],
  );

  static const LinearGradient accentGradient = LinearGradient(
    begin: Alignment.centerLeft,
    end: Alignment.centerRight,
    colors: [
      Color(0xFFFF3B30),
      Color(0xFFFF9500),
    ],
  );

  // Glassmorphism Tokens
  static const Color glassBase = Color(0xE6FFFFFF); // 90% opacity for better contrast
  static const Color glassBorder = Color(0x33FFFFFF);

  static ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      colorScheme: ColorScheme.fromSeed(
        seedColor: accent,
        primary: accent,
        surface: surface,
        onSurface: primaryText,
        secondary: Color(0xFF007AFF), // iOS-style blue
      ),
      scaffoldBackgroundColor: background,
      textTheme: GoogleFonts.outfitTextTheme().copyWith(
        titleLarge: GoogleFonts.outfit(
          color: primaryText,
          fontWeight: FontWeight.w700,
          fontSize: 32,
          letterSpacing: -0.8,
        ),
        bodyMedium: GoogleFonts.inter(
          color: secondaryText,
          fontSize: 16,
          height: 1.5,
        ),
      ),
      cardTheme: CardTheme(
        color: surface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: const BorderSide(color: border, width: 0.5),
        ),
      ),
    );
  }
}
