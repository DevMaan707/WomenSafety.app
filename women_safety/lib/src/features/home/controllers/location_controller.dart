import 'dart:async';

import 'package:get/get.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
import 'package:geolocator/geolocator.dart';

class LocationController extends GetxController {
  final Rx<LatLng?> currentLocation = Rx<LatLng?>(null);
  final RxBool isLoading = false.obs;
  final RxDouble zoomLevel = 18.0.obs;
  GoogleMapController? mapController;
  StreamSubscription<Position>? _positionStreamSubscription;

  @override
  void onInit() {
    super.onInit();
    _initializeLocation();
  }

  @override
  void onClose() {
    _positionStreamSubscription?.cancel();
    mapController?.dispose();
    super.onClose();
  }

  Future<void> _initializeLocation() async {
    final hasPermission = await _checkAndRequestPermission();
    if (hasPermission) {
      await getCurrentLocation();
    }
  }

  Future<bool> _checkAndRequestPermission() async {
    bool serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) {
      Get.snackbar(
        'Location Services Disabled',
        'Please enable location services in your device settings.',
        snackPosition: SnackPosition.BOTTOM,
      );
      return false;
    }

    LocationPermission permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        Get.snackbar(
          'Permission Denied',
          'Location permissions are required for this feature.',
          snackPosition: SnackPosition.BOTTOM,
        );
        return false;
      }
    }

    if (permission == LocationPermission.deniedForever) {
      Get.snackbar(
        'Permission Denied',
        'Location permissions are permanently denied. Please enable them in settings.',
        snackPosition: SnackPosition.BOTTOM,
      );
      return false;
    }

    return true;
  }

  Future<void> getCurrentLocation() async {
    if (isLoading.value) return;

    isLoading.value = true;
    try {
      final Position position = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.bestForNavigation,
        timeLimit: const Duration(seconds: 10),
      );

      currentLocation.value = LatLng(position.latitude, position.longitude);
      _animateToCurrentLocation();
      _startLocationStream();
    } catch (e) {
      Get.snackbar(
        'Error',
        'Failed to get location: ${e.toString()}',
        snackPosition: SnackPosition.BOTTOM,
      );
    } finally {
      isLoading.value = false;
    }
  }

  void _startLocationStream() {
    _positionStreamSubscription?.cancel();
    _positionStreamSubscription = Geolocator.getPositionStream(
      locationSettings: const LocationSettings(
        accuracy: LocationAccuracy.bestForNavigation,
        distanceFilter: 1,
      ),
    ).listen(
      (Position newPosition) {
        currentLocation.value =
            LatLng(newPosition.latitude, newPosition.longitude);
        _animateToCurrentLocation();
      },
      onError: (error) {
        Get.snackbar(
          'Error',
          'Location stream error: ${error.toString()}',
          snackPosition: SnackPosition.BOTTOM,
        );
      },
    );
  }

  void _animateToCurrentLocation() {
    if (mapController != null && currentLocation.value != null) {
      mapController!.animateCamera(
        CameraUpdate.newCameraPosition(
          CameraPosition(
            target: currentLocation.value!,
            zoom: zoomLevel.value,
          ),
        ),
      );
    }
  }

  void updateZoom(double newZoom) {
    if (newZoom < 1 || newZoom > 20) return; // Validate zoom level
    zoomLevel.value = newZoom;
    _animateToCurrentLocation();
  }

  void setMapController(GoogleMapController controller) {
    mapController = controller;
    _animateToCurrentLocation();
  }
}
