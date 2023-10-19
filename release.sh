#!/bin/sh

user="Travis CI"
email="travis@travis-ci.org"
perun="perun-linux-amd64"
github="https://$2@github.com/Appliscale"
release="https://github.com/Appliscale/perun/releases/download/$1/$perun.tar.gz"
files="https://raw.githubusercontent.com/Appliscale/perun/master"

cd ~
sudo apt-get install rpm
git clone https://github.com/Appliscale/rpmbuild.git
cd rpmbuild/SOURCES 
rm $perun.tar.gz
wget $release
wget $files/LICENSE
tar xvzf $perun.tar.gz
rm $perun.tar.gz
tar cvzf $perun.tar.gz $perun LICENSE
rm LICENSE $perun
cd ..
rpmbuild -ba SPECS/$perun.spec
git remote
git config user.email $email
git config user.name $user
git add .
git commit -m "[AUTO] Update RPM by Travis CI. Perun $1"
git push $github/rpmbuild.git master

cd ~
git clone https://github.com/Appliscale/perun-dpkg.git
cd perun-dpkg
chmod +x control.sh
./control.sh $1
cd /perun/usr/local/bin
rm $perun
wget $release
tar xvzf $perun.tar.gz
rm $perun.tar.gz
cd ~/perun-dpkg
dpkg-deb --build perun
git remote
git config user.email $email
git config user.name $user
git add .
git commit -m "[AUTO] Update DPKG by Travis CI. Perun $1"
git push $github/perun-dpkg.git master

cd ~
git clone https://github.com/Appliscale/homebrew-tap.git
cd homebrew-tap
wget https://github.com/Appliscale/perun/blob/master/formula.sh
chmod +x formula.sh
./formula.sh $1
rm formula.sh
git remote
git config user.email $email
git config user.name $user
git add .
git commit -m "[AUTO] Update Homebrew by Travis CI. Perun $1"
git push $github/homebrew-tap.git master
