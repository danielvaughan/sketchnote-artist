import 'package:flutter/material.dart';
import 'theme/app_theme.dart';
import 'screens/home_page.dart';

void main() {
  runApp(const SketchnoteApp());
}

class SketchnoteApp extends StatelessWidget {
  const SketchnoteApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Sketchnote Artist',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      home: const HomePage(),
    );
  }
}
