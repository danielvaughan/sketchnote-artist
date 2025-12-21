class AgentEvent {
  final Map<String, dynamic>? content;
  final List<dynamic>? models;
  final String? error;

  AgentEvent({this.content, this.models, this.error});

  factory AgentEvent.fromJson(Map<String, dynamic> json) {
    return AgentEvent(
      content: json['content'] as Map<String, dynamic>?,
      models: json['models'] as List<dynamic>?,
      error: json['error'] as String?,
    );
  }
}
