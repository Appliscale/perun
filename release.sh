#!/bin/sh

user="Travis CI"
email="travis@travis-ci.org"
perun="perun-linux-amd64"
github="https://$2@github.com/Appliscale"
release="https://github.com/Appliscale/perun/releases/download/$1/$perun.tar.gz"
files="https://raw.githubusercontent.com/Appliscale/perun/master"

sudo apt-get install rpm
git clone https://github.com/Appliscale/rpmbuild.git
cd rpmbuild/SOURCES 
rm $perun.tar.gz
wget $release
wget $files/defaults/main.yaml
wget $files/LICENSE
tar xvzf $perun.tar.gz
rm $perun.tar.gz
tar cvzf $perun.tar.gz $perun LICENSE main.yaml
rm LICENSE main.yaml $perun
cd ..
rpmbuild -ba /SPEC/$perun.spec
git remote
git config user.email $email
git config user.name $user
git add .
git commit -m "[AUTO] Update RPM by Travis CI"
git push $github/rpmbuild.git master

cd ~
git clone https://github.com/Appliscale/perun-dpkg.git
cd perun-dpkg/perun/usr/local/bin
rm $perun
wget $release
tar xvzf $perun.tar.gz
rm $perun.tar.gz
cd ~
cd perun-dpkg
dpkg-deb --build perun
git remote
git config user.email $email
git config user.name $user
git add .
git commit -m "[AUTO] Update DPKG by Travis CI"
git push $github/perun-dpkg.git master
