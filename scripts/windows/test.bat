XCOPY /Y /I build wintest
TYPE settings.toml |sed -e s/hogehoge@/%email%/g -e s/hogehoge/%password%/g>.\wintest\settings.toml
cd .\wintest
eliminateToil.exe nikkei
