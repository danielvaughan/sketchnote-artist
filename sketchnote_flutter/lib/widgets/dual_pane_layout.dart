import 'package:flutter/material.dart';

class DualPaneLayout extends StatelessWidget {
  final Widget left;
  final Widget right;

  const DualPaneLayout({
    super.key,
    required this.left,
    required this.right,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        if (constraints.maxWidth > 900) {
          return Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(child: left),
              const SizedBox(width: 32),
              Expanded(child: right),
            ],
          );
        } else {
          return Column(
            children: [
              left,
              const SizedBox(height: 32),
              right,
            ],
          );
        }
      },
    );
  }
}
