**SAGE** Status

Back in the day when I was a young-un, I remember the following things:
    {{ range . }}✅ `{{ .data.ID }}`
            Latest block: `{{ .latestBlock }}`
            RPC: `{{ .data.RPC }}`
            LCD: `{{ .data.LCD }}`
            Catching up: `{{ .sync.SyncInfo.CatchingUp }}`
    {{ end }}