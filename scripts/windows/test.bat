XCOPY /Y /I build wintest
TYPE settings.toml |sed -e s/hogehoge@/%email%/g -e s/hogehoge/%password%/g -e "s/  proxy/# proxy/g">.\wintest\settings.toml
cd .\wintest
eliminateToil.exe nikkei
