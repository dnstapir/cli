Name:          dnstapir-cli
# NOTE: Version must match VERSION file - validated by Makefile srpm target
Version:       v0.3
Release:       1%{?dist}
Group:         dnstapir/edge
Summary:       DNS TAPIR EDGE Cli Tool
License:       BSD
URL:           https://www.github.com/dnstapir/cli
Source0:       %{name}-%{version}.tar.gz
Source1:       dnstapir-renew.service
Source2:       dnstapir-renew.timer
BuildRequires: git
BuildRequires: golang

%description
DNS TAPIR EDGE ClI Tool for managing an EDGE deployment

# Disable building of debug packages for RHEL (we include symbols per default)
%if 0%{?rhel} >= 9
%global debug_package %{nil}
%endif

%{!?_unitdir: %define _unitdir /usr/lib/systemd/system/}
%{!?_sysusersdir: %define _sysusersdir /usr/lib/sysusers.d/}

%prep
%setup -n %{name}

%build
make

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_unitdir}
mkdir -p %{buildroot}%{_sysconfdir}/dnstapir/certs

install -p -m 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -m 0644 %{SOURCE1} %{buildroot}%{_unitdir}
install -m 0644 %{SOURCE2} %{buildroot}%{_unitdir}

%files
%attr(0770,root,dnstapir) %dir %{_sysconfdir}/dnstapir
%attr(0770,root,dnstapir) %dir %{_sysconfdir}/dnstapir/certs
%attr(0770,root,dnstapir) %{_bindir}/%{name}
%attr(0660,-,dnstapir) %ghost %{_sysconfdir}/dnstapir/dnstapir-cli.yaml
%attr(0644,root,dnstapir) %{_unitdir}/dnstapir-renew.service
%attr(0644,root,dnstapir) %{_unitdir}/dnstapir-renew.timer

%pre
/usr/bin/getent group dnstapir || /usr/sbin/groupadd -r dnstapir
/usr/bin/getent passwd dnstapir-renew || /usr/sbin/useradd -r -d /etc/dnstapir -G dnstapir -s /sbin/nologin dnstapir-renew

%post

%preun

%postun

%check

%changelog
