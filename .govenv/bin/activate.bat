@echo off

if defined GOVENV_ENABLE (
echo Virtual environment has already been activated.
GOTO fin
) else (
GOTO setvariables
)

:setvariables
SET "GOVENV_OLD_PATH=%PATH%"
SET "GOVENV_OLD_PROMPT=%PROMPT%"

SET GOVENV_ENABLE=1
SET "GOVENV_PROJECT=slide_converter"
SET "GOVENV_MANAGEMENT_DIR=C:\Users\nakagash\.govenv"
SET "GOVENV_DATA_DIR=.govenv"

SET "GOROOT=C:\Users\nakagash\.govenv\goroots\go1.9.2"
SET "GOPATH=C:\Users\nakagash\.govenv\tools\bin\slide_converter"

SET "PROMPT=(slide_converter)%PROMPT%"
SET "PATH=%GOPATH%\.govenv\bin;%GOROOT%\bin;%PATH%"

:fin
