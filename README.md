# texApi

Installation is easy

Put this repo folder in:
**~/tex_backend/texApi**

Run initial configuration:
```bash
make init-sys
```
After that carefully configure the **systemd/system/texApi.service** file

Update the app (need to have github access keys):
```bash
bash ~/tex_backend/texApi/scripts/update_tex.sh
```

Write absolute path of uploads directory to .env, then run:
```bash
make upload-dir
```

```bash
make db
make dev
```

To build application:

```bash
make build
```


