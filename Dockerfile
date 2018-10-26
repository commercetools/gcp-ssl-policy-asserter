#Scratch container with CA certs from alpine linux
FROM drone/ca-certs
ADD policy_asserter /
CMD ["/policy_asserter"]