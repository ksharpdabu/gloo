// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/setup"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Needed to run tests in GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// From https://github.com/kubernetes/client-go/blob/53c7adfd0294caa142d961e1f780f74081d5b15f/examples/out-of-cluster-client-configuration/main.go#L31
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("V1Emitter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace1        string
		namespace2        string
		cfg               *rest.Config
		emitter           StatusEmitter
		kubeServiceClient KubeServiceClient
		ingressClient     IngressClient
	)

	BeforeEach(func() {
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		err := setup.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace2)
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		var kube kubernetes.Interface
		// KubeService Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		kubeServiceClientFactory := &factory.KubeConfigMapClientFactory{
			Clientset: kube,
		}
		kubeServiceClient, err = NewKubeServiceClient(kubeServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// Ingress Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		ingressClientFactory := &factory.KubeConfigMapClientFactory{
			Clientset: kube,
		}
		ingressClient, err = NewIngressClient(ingressClientFactory)
		Expect(err).NotTo(HaveOccurred())
		emitter = NewStatusEmitter(kubeServiceClient, ingressClient)
	})
	AfterEach(func() {
		setup.TeardownKube(namespace1)
		setup.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *StatusSnapshot

		/*
			KubeService
		*/

		assertSnapshotServices := func(expectServices KubeServiceList, unexpectServices KubeServiceList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectServices {
						if _, err := snap.Services.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectServices {
						if _, err := snap.Services.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := kubeServiceClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := kubeServiceClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}

		kubeService1a, err := kubeServiceClient.Write(NewKubeService(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		kubeService1b, err := kubeServiceClient.Write(NewKubeService(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServices(KubeServiceList{kubeService1a, kubeService1b}, nil)

		kubeService2a, err := kubeServiceClient.Write(NewKubeService(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		kubeService2b, err := kubeServiceClient.Write(NewKubeService(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServices(KubeServiceList{kubeService1a, kubeService1b, kubeService2a, kubeService2b}, nil)

		err = kubeServiceClient.Delete(kubeService2a.Metadata.Namespace, kubeService2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = kubeServiceClient.Delete(kubeService2b.Metadata.Namespace, kubeService2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServices(KubeServiceList{kubeService1a, kubeService1b}, KubeServiceList{kubeService2a, kubeService2b})

		err = kubeServiceClient.Delete(kubeService1a.Metadata.Namespace, kubeService1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = kubeServiceClient.Delete(kubeService1b.Metadata.Namespace, kubeService1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotServices(nil, KubeServiceList{kubeService1a, kubeService1b, kubeService2a, kubeService2b})

		/*
			Ingress
		*/

		assertSnapshotIngresses := func(expectIngresses IngressList, unexpectIngresses IngressList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectIngresses {
						if _, err := snap.Ingresses.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectIngresses {
						if _, err := snap.Ingresses.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := ingressClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := ingressClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}

		ingress1a, err := ingressClient.Write(NewIngress(namespace1, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		ingress1b, err := ingressClient.Write(NewIngress(namespace2, "angela"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b}, nil)

		ingress2a, err := ingressClient.Write(NewIngress(namespace1, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		ingress2b, err := ingressClient.Write(NewIngress(namespace2, "bob"), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b, ingress2a, ingress2b}, nil)

		err = ingressClient.Delete(ingress2a.Metadata.Namespace, ingress2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = ingressClient.Delete(ingress2b.Metadata.Namespace, ingress2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(IngressList{ingress1a, ingress1b}, IngressList{ingress2a, ingress2b})

		err = ingressClient.Delete(ingress1a.Metadata.Namespace, ingress1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = ingressClient.Delete(ingress1b.Metadata.Namespace, ingress1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIngresses(nil, IngressList{ingress1a, ingress1b, ingress2a, ingress2b})
	})
})
