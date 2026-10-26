package plugin_test

import (
	"time"

	"github.com/meln5674/grafana-mongodb-community-plugin/pkg/plugin"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ApplyMacros", func() {
	from := time.UnixMilli(1700000000000)
	to := time.UnixMilli(1700003600000)
	interval := 15 * time.Second

	It("Should replace both macro forms", func() {
		Expect(plugin.ApplyMacros(
			`{"$gte": {"$date": $__from}, "$lte": {"$date": ${__to}}, "step": $__interval_ms, "s2": ${__interval_ms}}`,
			from, to, interval,
		)).To(Equal(`{"$gte": {"$date": 1700000000000}, "$lte": {"$date": 1700003600000}, "step": 15000, "s2": 15000}`))
	})

	It("Should leave unknown variables and longer names alone", func() {
		Expect(plugin.ApplyMacros(`$__fromage $__range $myvar`, from, to, interval)).
			To(Equal(`$__fromage $__range $myvar`))
	})
})

var _ = Describe("ConvertGoTimeFormatToMongo", func() {
	It("Should convert a truncated ruby date", func() {
		converted, err := plugin.ConvertGoTimeFormatToMongo(time.RubyDate[4:])
		Expect(err).ToNot(HaveOccurred())
		Expect(converted).To(Equal("%b %d %H:%M:%S %z %Y"))
	})
})
