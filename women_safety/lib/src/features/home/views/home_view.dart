import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import '../../../components/report_dialog.dart';
import '../controllers/location_controller.dart';
import '../controllers/audio_controller.dart';

class MainView extends GetView<LocationController> {
  const MainView({super.key});

  @override
  Widget build(BuildContext context) {
    final audioController = Get.put(AudioController());

    return Scaffold(
      body: Stack(
        children: [
          Obx(() => GoogleMap(
                initialCameraPosition: CameraPosition(
                  target:
                      controller.currentLocation.value ?? const LatLng(0, 0),
                  zoom: controller.zoomLevel.value,
                ),
                myLocationEnabled: true,
                myLocationButtonEnabled: false,
                compassEnabled: true,
                zoomControlsEnabled: false,
                mapType: MapType.normal,
                onMapCreated: (GoogleMapController mapController) {
                  controller.setMapController(mapController);
                  controller.mapController = mapController;
                },
              )),
          Obx(() => controller.isLoading.value
              ? const Center(child: CircularProgressIndicator())
              : const SizedBox.shrink()),
          Positioned(
            top: 50,
            right: 16,
            child: FloatingActionButton(
              mini: true,
              backgroundColor: Colors.white,
              child: const Icon(Icons.my_location, color: Colors.black87),
              onPressed: controller.getCurrentLocation,
            ),
          ),
          Positioned(
            bottom: 32,
            left: 0,
            right: 0,
            child: Center(
              child: SafetyButton(
                onPressed: () async {
                  await audioController.toggleRecording();
                  Get.snackbar(
                    'Safety Feature',
                    'Current Location: ${controller.currentLocation.value}',
                    backgroundColor: Colors.white,
                  );
                },
              ),
            ),
          ),
          Positioned(
            top: 50,
            left: 16,
            child: FloatingActionButton(
              mini: true,
              backgroundColor: Colors.white,
              onPressed: () {
                Get.dialog(
                  ReportIncidentDialog(),
                );
              },
              child: const Icon(
                Icons.report_problem_outlined,
                color: Color(0xFFFF4081),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class SafetyButton extends StatelessWidget {
  final VoidCallback onPressed;
  const SafetyButton({Key? key, required this.onPressed}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GetX<AudioController>(
      builder: (audioController) => Material(
        elevation: 8,
        borderRadius: BorderRadius.circular(28),
        color: audioController.isRecording.value
            ? Colors.red
            : const Color(0xFFFF4081),
        child: InkWell(
          onTap: onPressed,
          borderRadius: BorderRadius.circular(28),
          child: Container(
            padding: const EdgeInsets.symmetric(
              horizontal: 24,
              vertical: 16,
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  audioController.isRecording.value ? Icons.stop : Icons.shield,
                  color: Colors.white,
                  size: 24,
                ),
                const SizedBox(width: 8),
                Text(
                  audioController.isRecording.value
                      ? 'Stop Recording'
                      : 'Keep Me Safe',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
