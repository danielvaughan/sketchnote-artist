import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/agent_event.dart';

class SseService {
  final String baseUrl;
  http.Client? _client;

  SseService({required this.baseUrl});

  Stream<AgentEvent> runSse({
    required String userId,
    required String sessionId,
    required String message,
  }) async* {
    final url = '$baseUrl/run_sse';
    _client = http.Client();

    final body = {
      'appName': 'sketchnote-artist',
      'userId': userId,
      'sessionId': sessionId,
      'newMessage': {
        'role': 'user',
        'parts': [{'text': message}]
      },
      'streaming': true
    };

    final request = http.Request('POST', Uri.parse(url));
    request.headers.addAll({
      'Content-Type': 'application/json',
      'Accept': 'text/event-stream',
    });
    request.body = jsonEncode(body);

    try {
      final response = await _client!.send(request);
      
      if (response.statusCode != 200) {
        yield AgentEvent(error: 'Failed to connect: ${response.statusCode}');
        return;
      }

      // Read the stream line by line
      String buffer = '';
      await for (final chunk in response.stream.transform(utf8.decoder)) {
        buffer += chunk;
        final lines = buffer.split('\n\n');
        buffer = lines.removeAt(lines.length - 1);

        for (final line in lines) {
          if (line.isEmpty) continue;
          if (line.startsWith(':')) continue; // Heartbeats

          if (line.startsWith('event: error')) {
            yield AgentEvent(error: 'Stream error');
            continue;
          }

          if (line.startsWith('data: ')) {
            final dataStr = line.substring(6).trim();
            if (dataStr.isEmpty) continue;
            try {
              yield AgentEvent.fromJson(jsonDecode(dataStr));
            } catch (e) {
              // Ignore partial or malformed data
            }
          }
        }
      }
    } catch (e) {
      yield AgentEvent(error: 'Network error: $e');
    }
  }

  void stopSse() {
    _client?.close();
    _client = null;
  }
}
