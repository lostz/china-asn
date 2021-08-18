# china-asn
中国区所有的asn通过https://whois.ipip.net/countries/CN 获得
默认把所有的中国asn号存储在/etc/cn.conf 里
推荐使用systemd.timer 和systemd.path 来实现bird 自动更新
```
go build -o cnasn main.go

# /etc/systemd/system/cnasn-gen.servcie
[Unit]
Description=china asn

[Service]
Type=oneshot
ExecStart=/usr/local/bin/cnasn
# /etc/systemd/system/cnasn-gen.timer
[Unit]
Description=Runs  every day

[Timer]
OnUnitActiveSec=1d
Unit=cnasn-gen.service

[Install]
WantedBy=multi-user.target

# /etc/systemd/system/cnasn.path
[Unit]
Description=Watch /etc/cn.conf changes

[Path]
PathModified=/etc/cn.conf

[Install]
WantedBy=multi-user.target


# /etc/systemd/system/cnasn.servcie
[Unit]
Description=Reload bird
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl reload bird

[Install]
WantedBy=multi-user.target

```



