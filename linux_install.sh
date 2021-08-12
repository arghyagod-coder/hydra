#!/bin/bash

echo "Downloading hydra..."
wget https://github.com/Shravan-1908/hydra/releases/latest/download/hydra-linux-amd64 

echo "Adding hydra into PATH..."

mkdir -p ~/.hydra

mv ./hydra-linux-amd64 ~/.hydra/hydra
chmod +x ~/.hydra/hydra
echo "export PATH=$PATH:~/.hydra" >> ~/.bashrc
fish -c "set -U fish_user_paths ~/.hydra/ $fish_user_paths"      
echo "export PATH=$PATH:~/.hydra" >> ~/.zshrc

echo "hydra installation is completed!"
echo "You need to restart the shell to use hydra."
