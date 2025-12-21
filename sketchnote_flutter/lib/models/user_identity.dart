class UserIdentity {
  final String email;

  UserIdentity({required this.email});

  factory UserIdentity.fromJson(Map<String, dynamic> json) {
    return UserIdentity(
      email: json['email'] as String? ?? 'local-user',
    );
  }
}
