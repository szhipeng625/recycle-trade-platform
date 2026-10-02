@echo off
chcp 65001 >nul
title 废旧物资回收交易平台 - 启动脚本

echo ========================================
echo   废旧物资回收交易平台
echo   前端: Vue3 + Element Plus + ECharts
echo   后端: Go + Gin + GORM
echo   数据库: MySQL 8.0
echo ========================================
echo.

REM 检查 MySQL 是否运行
tasklist /fi "imagename eq mysqld.exe" 2>nul | findstr mysqld >nul
if %errorlevel% neq 0 (
    echo [1/3] 启动 MySQL...
    start "MySQL" /min "D:\mysql-8.0.40-winx64\bin\mysqld.exe" --console
    timeout /t 5 /nobreak >nul
) else (
    echo [1/3] MySQL 已在运行
)

echo [2/3] 启动后端服务 (端口 9000)...
cd /d "%~dp0server"
start "后端服务-Go" cmd /k "server.exe"

timeout /t 2 /nobreak >nul

echo [3/3] 启动前端服务 (端口 5173)...
cd /d "%~dp0client"
start "前端服务-Vue" cmd /k "npx vite --host"

timeout /t 3 /nobreak >nul
echo.
echo ========================================
echo   启动完成！
echo   前端地址: http://localhost:5173
echo   后端地址: http://localhost:9000
echo   管理员账号: admin / admin123
echo ========================================
echo.
pause
