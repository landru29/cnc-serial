package main

import (
	"go.bug.st/serial"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"io"
	"fmt"
	"context"
	"time"
	"errors"
)


var (
	lastCommand []string
	cursor int
)

func sendCommand(port io.ReadWriter, bottomArea io.Writer, text string) {
	fmt.Fprintf(bottomArea, " > %s\n", text)

	if _, err := fmt.Fprintf(port, "%s\n", text); err != nil {
		fmt.Fprintf(bottomArea, " - ERR %s\n", err.Error())
	}

	lastCommand = append(lastCommand, text)
	cursor = len(lastCommand)-1
}


func mainCommand() *cobra.Command {
	var (
		portName        string
		defaultPortName string
		bitRate         int
	)

	ports, err := serial.GetPortsList()
	if err==nil && len(ports)>0 {
		defaultPortName=ports[0]
	}

	output := &cobra.Command{
		Use: "serial",
		Short: "Serial monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := serial.Open(portName, &serial.Mode{
				BaudRate: bitRate,
			})
			if err != nil {
				return err
			}

			defer func() {
				_ = port.Close()
			}()

			app := tview.NewApplication()

			bottomArea := tview.NewTextView().
				SetDynamicColors(true).
				SetRegions(true).
				SetChangedFunc(func() {
					app.Draw()
				})

			go func() {
				for {
					buf := make([]byte, 200)

					n, err := port.Read(buf)

					switch {
					case errors.Is(err, io.EOF):
						// Do nothing
					case err != nil:
						fmt.Fprintf(bottomArea, "> ERR %s\n", err.Error())
					default:
						fmt.Fprintf(bottomArea, "  < %s", string(buf[:n]))
					}

					time.Sleep(2000*time.Millisecond)
				}
			}()

			var input *tview.InputField

			input = tview.NewInputField().
				SetLabel("Enter command").
				SetFieldWidth(80).
				SetPlaceholder("here").
				SetDoneFunc(func(key tcell.Key){
					text := input.GetText()

					if text == "exit" {
						app.Stop()

						return
					}

					sendCommand(port, bottomArea, text)

					input.SetText("")
				})
			input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
					if event.Key() == tcell.KeyUp && cursor>=0 {
						input.SetText(lastCommand[cursor])
						cursor--
					}

					if event.Key() == tcell.KeyDown && cursor>=0 && cursor<len(lastCommand)-1 {
						input.SetText(lastCommand[cursor])
						cursor++
					}

					return event
				})

			bottomArea.SetBorder(true)

			flex := tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(bottomArea, 0, 6, false).
				AddItem(input, 0, 1, true)


			if err := app.SetRoot(flex, true).Run(); err != nil {
				return err
			}

			return nil
		},
	}

	output.Flags().IntVarP(&bitRate, "bit-rate", "b", 115200, "Bit rate")
	output.Flags().StringVarP(&portName, "port", "p", defaultPortName, "Port name")

	return output
}

func main() {
	if err := mainCommand().ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
