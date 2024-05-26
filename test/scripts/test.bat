@echo off
echo "BAT"
echo FLAGS: %*

rem Використання циклу для виведення кожного параметру по черзі
for %%G in (%*) do (
    echo %%G
)