import 'package:cc_essentials/cc_essentials.dart';
import 'package:cc_essentials/helpers/logging/logger.dart';
import 'package:cc_essentials/services/shared_preferences/shared_preference_service.dart';
import 'package:cc_essentials/theme/custom_theme.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'src/features/home/bindings/home_bindings.dart';
import 'src/features/home/views/home_view.dart';

final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();
Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await CCEssentials.initialize(
      primaryColor: Colors.orange,
      accentColor: Colors.teal,
      navigatorKey: navigatorKey);
  logger.i(SharedPreferencesService().isLoggedIn());
  runApp(Sathee());
}

class Sathee extends StatelessWidget {
  const Sathee({super.key});
  @override
  Widget build(BuildContext context) {
    return GetMaterialApp(
      theme: CustomTheme.lightTheme(),
      darkTheme: CustomTheme.darkTheme(),
      themeMode: ThemeMode.system,
      initialBinding: HomeBinding(),
      home: MainView(),
    );
  }
}
