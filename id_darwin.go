// +build darwin

package machineid

//#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
//#include <CoreFoundation/CoreFoundation.h>
//#include <IOKit/IOKitLib.h>
//
//const char * getHardwareUUID() {
//    io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault,
//		IOServiceMatching("IOPlatformExpertDevice"));
//    if (service == 0) {
//        return NULL;
//    }
//
//    CFStringRef uuid = IORegistryEntryCreateCFProperty(service, CFSTR("IOPlatformUUID"), kCFAllocatorDefault, 0);
//    if (uuid != NULL) {
//        IOObjectRelease(service);
//        return CFStringGetCStringPtr(uuid, kCFStringEncodingUTF8);
//    }
//
//    IOObjectRelease(service);
//
//    return NULL;
//}
import "C"
import "errors"

// machineID returns the uuid stored in IO Registry.
// If there is an error an empty string is returned.
func machineID() (string, error) {
	uuid := C.GoString(C.getHardwareUUID())
	if uuid == "" {
		return "", errors.New("failed to retrieve the property in registry")
	}
	return uuid, nil
}