// +build linux

package machineid

import (
	"log"
	"testing"
)

func TestIsContainer(t *testing.T) {
	wanted := isContainer()
	if wanted {
		log.Println("Inside the container")
	} else {
		log.Println("Currently not in container")
	}
}

func TestExtractMountInfoId(t *testing.T) {
	mount := []string {
		"767 553 0:187 / / rw,relatime master:167 - overlay overlay rw,lowerdir=/var/lib/docker/overlay2/l/ZTHS22AHCPD2HJEY6UIKIO3BHY:/var/lib/docker/overlay2/l/4P6QELQ6532G5362S5VVTA7Y7K,upperdir=/var/lib/docker/overlay2/76c8877e95fa589df1fb97bf831ec221df130fdfb8f1f1cb8166bd99bebf51de/diff,workdir=/var/lib/docker/overlay2/76c8877e95fa589df1fb97bf831ec221df130fdfb8f1f1cb8166bd99bebf51de/work",
	}
	for _, m := range mount {
		id := extractMountInfoId(m)
		if id == "" {
			log.Printf("Failed to extract the id from %s", m)
		}
		log.Printf("Extract id %s from %s", id, m)
	}
}

func TestExtractCtrlGroupId(t *testing.T) {
	cgroups := []string {
		"1:cpuset:/\n1:cpuset:/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-podf8e932ad_5487_4514_a5be_b75ad1b7a6ce.slice/crio-ee55f03bd921c55955d8995a0adbb9f19352603a637ea27f6ca8397b715435eb.scope",
		"5:net_cls:/system.slice/docker-afd862d2ed48ef5dc0ce8f1863e4475894e331098c9a512789233ca9ca06fc62.scope\n1:cpuset:/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-podf8e932ad_5487_4514_a5be_b75ad1b7a6ce.slice",
		"5:cpuset:/\n5:net_prio,net_cls:/docker/de630f22746b9c06c412858f26ca286c6cdfed086d3b302998aa403d9dcedc42\n5:net_prio,net_cls:/",
		"3:net_cls:/kubepods/burstable/pod5f399c1a-f9fc-11e8-bf65-246e9659ebfc/9170559b8aadd07d99978d9460cf8d1c71552f3c64fefc7e9906ab3fb7e18f69",
	}

	for _, cgroup := range cgroups {
		id := extractCtrlGroupId(cgroup)
		if id == "" {
			log.Printf("Failed to extract the id from %s", cgroup)
		}
		log.Printf("Extract id %s from %s", id, cgroup)
	}
}