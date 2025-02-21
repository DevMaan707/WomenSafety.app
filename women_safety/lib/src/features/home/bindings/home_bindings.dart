import 'package:get/get.dart';
import '../controllers/location_controller.dart';

class HomeBinding extends Bindings {
  @override
  void dependencies() {
    Get.put(LocationController());
  }
}
