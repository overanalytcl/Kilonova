package checkers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/KiloProjects/kilonova/eval"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func legacyCheckerTask(ctx context.Context, box eval.Sandbox, job *customCheckerInput) (*checkerResult, error) {
	rez := &checkerResult{}
	lang, ok := eval.Langs[eval.GetLangByFilename(job.c.filename)]
	if !ok {
		rez.Output = ErrOut
		return rez, nil
	}

	if ok := copyFiles(box, lang, job); !ok {
		rez.Output = ErrOut
		return rez, nil
	}

	goodCmd, err := eval.MakeGoodCommand(lang.RunCommand)
	if err != nil {
		rez.Output = ErrOut
		return rez, nil
	}
	goodCmd = append(goodCmd, "/box/program.out", "/box/correct.out", "/box/correct.in")

	conf := &eval.RunConfig{
		OutputPath: "/box/checker_verdict.out",

		MemoryLimit: checkerMemoryLimit,

		WallTimeLimit: 20,
	}

	if _, err := box.RunCommand(ctx, goodCmd, conf); err != nil {
		rez.Output = ErrOut
		return rez, nil
	}

	var out bytes.Buffer
	if err := box.ReadFile("/box/checker_verdict.out", &out); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			zap.S().Warn("Couldn't read checker output: ", err)
		}
		out.Reset()
	}

	var percentage float64
	if _, err := fmt.Fscanf(&out, "%f ", &percentage); err != nil {
		rez.Output = "Wrong checker output"
		return rez, nil
	}
	rez.Percentage = decimal.NewFromFloat(percentage)

	rez.Output = strings.TrimSpace(out.String())
	return rez, nil
}
