@echo off
REM Docker Cleanup Script for Windows
REM This script removes unused Docker containers, images, networks, and volumes

setlocal enabledelayedexpansion

echo Starting Docker cleanup process...

echo Current Docker resources:
docker ps -aq 2>nul | find /c /v "" > temp_count.txt
set /p containers=<temp_count.txt
docker images -q 2>nul | find /c /v "" > temp_count.txt
set /p images=<temp_count.txt
docker volume ls -q 2>nul | find /c /v "" > temp_count.txt
set /p volumes=<temp_count.txt
docker network ls -q 2>nul | find /c /v "" > temp_count.txt
set /p networks=<temp_count.txt
del temp_count.txt

echo   Containers: %containers%
echo   Images: %images%
echo   Volumes: %volumes%
echo   Networks: %networks%

echo.

REM Remove stopped containers
echo Removing stopped containers...
for /f %%i in ('docker ps -aq --filter "status=exited" --filter "status=created"') do (
    docker rm %%i
    if !errorlevel! equ 0 (
        set /a "stopped_count+=1"
    )
)

REM Remove dangling images
echo Removing dangling images...
for /f %%i in ('docker images -f "dangling=true" -q') do (
    docker rmi %%i
    if !errorlevel! equ 0 (
        set /a "dangling_count+=1"
    )
)

REM Clean system (removes unused data)
echo Cleaning Docker system...
docker system prune -f

REM Show disk usage after cleanup
echo.
echo Docker disk usage after cleanup:
docker system df

echo.
echo Cleanup completed!
if defined stopped_count echo   Removed %stopped_count% stopped containers
if defined dangling_count echo   Removed %dangling_count% dangling images

echo.
echo Final Docker resources:
docker ps -aq 2>nul | find /c /v "" > temp_count.txt
set /p containers=<temp_count.txt
docker images -q 2>nul | find /c /v "" > temp_count.txt
set /p images=<temp_count.txt
docker volume ls -q 2>nul | find /c /v "" > temp_count.txt
set /p volumes=<temp_count.txt
docker network ls -q 2>nul | find /c /v "" > temp_count.txt
set /p networks=<temp_count.txt
del temp_count.txt

echo   Containers: %containers%
echo   Images: %images%
echo   Volumes: %volumes%
echo   Networks: %networks%

endlocal