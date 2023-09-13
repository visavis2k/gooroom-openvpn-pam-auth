# gooroom-openvpn-pam-auth

구름 사용자가 로그인할 때 입력한 인증 정보를 사용해 OpenVPN에 연결합니다.

## 적용 방법

1. 프로젝트를 빌드합니다.
2. 생성된 gooroom-openvpn-pam-auth 파일을 /usr/bin 디렉터리로 이동합니다.
3. /etc/pam.d/lightdm 파일에 다음 내용을 추가합니다.

```
auth optional pam_exec.so debug expose_authtok type=auth log=/var/log/gooroom-pam.log /usr/bin/gooroom-openvpn-pam-auth
```

4. 이 프로그램은 /etc/openvpn/client/gooroom.ovpn OpenVPN 클라이언트 구성 파일을 사용해 OpenVPN 서버에 연결합니다. OpenVPN 클라이언트 구성 파일의 예는 다음과 같습니다.

```
client
proto udp
explicit-exit-notify
remote 10.0.2.6 1194
connect-retry 0
connect-retry-max 0
dev tun
resolv-retry infinite
nobind
persist-key
persist-tun
verify-x509-name vpn.javaworld.co.kr name
auth SHA256
auth-nocache
cipher AES-128-GCM
tls-client
tls-version-min 1.2
tls-cipher TLS-ECDHE-ECDSA-WITH-AES-128-GCM-SHA256
ignore-unknown-option block-outside-dns
verb 3
setenv verb 9
<ca>
-----BEGIN CERTIFICATE-----
MIIDBTCCAe2gAwIBAgIUR3KuIJbeNwclvplihrmLmhmQFKAwDQYJKoZIhvcNAQEL
... 생략 ...
1QfiMGJ1+9FD
-----END CERTIFICATE-----
</ca>
<tls-crypt>
#
# 2048 bit OpenVPN static key
#
-----BEGIN OpenVPN Static key V1-----
d66a751c8713458664f297bba1a4ee4f
... 생략 ...
71387cd7975b0b76950a72996e4acaf5
-----END OpenVPN Static key V1-----
</tls-crypt>
```

pam_exec에 대해선 [링크](https://wariua.github.io/linux-pam-docs-ko/sag-pam_exec.html)를 참고하세요.

