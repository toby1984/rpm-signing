Name:           my-payload-package
Version:        1.0
Release:        1%{?dist}
Summary:        Minimal package to deploy payload.txt

License:        MIT
BuildArch:      noarch

%description
A minimal RPM package that installs payload.txt to the root directory.

%prep
# Nothing to unpack

%build
# Nothing to compile

%install
mkdir -p %{buildroot}/
cp %{_sourcedir}/payload.txt %{buildroot}/payload.txt

%files
%attr(0640, root, root) /payload.txt
