import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/user_identity.dart';

class ApiService {
  final String baseUrl;

  ApiService({required this.baseUrl});

  Future<UserIdentity> getUserIdentity() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/me'));
      if (response.statusCode == 200) {
        return UserIdentity.fromJson(jsonDecode(response.body));
      }
    } catch (e) {
      // Error fetching identity
    }
    return UserIdentity(email: 'local-user');
  }

  Future<String> createSession(String userId) async {
    const appName = 'sketchnote-artist';
    final response = await http.post(
      Uri.parse('$baseUrl/apps/$appName/users/$userId/sessions'),
    );
    if (response.statusCode == 200 || response.statusCode == 201) {
      final data = jsonDecode(response.body);
      return data['id'] as String;
    }
    throw Exception('Failed to create session: ${response.statusCode}');
  }
}
