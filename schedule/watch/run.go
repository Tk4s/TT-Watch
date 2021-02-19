package watch

import (
	"sync"

	"github.com/gogf/gf/os/gcron"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	w := NewWatch()
	if w == nil {
		return
	}

	defer w.close()

	influence := []string{"binance", "elonmusk"}

	gcron.AddSingleton("0 */1 * * * *", func() {
		wg := sync.WaitGroup{}
		wg.Add(len(influence))
		for _, star := range influence {
			go w.do(star, wg)
		}
		wg.Wait()

	})

	select {}
}
