**SAGE** Status

Back in the day when I was a young-un, I remember the following things:
    {{ range . }}✅ `{{ .ID }}`
            RPC: `{{ .RPC }}`
            LCD: `{{ .LCD }}`
            Latest block: `{{ .LatestBlock }}`
    {{ end }}