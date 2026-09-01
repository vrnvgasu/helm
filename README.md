# Создаём Helm-релиз в Kubernetes
1. запускаем проект в Kubernetes
   helm install echo-server ./echo-server
   
2. Проверяем, что релиз установлен
   helm list
             
3. Проверяем, что поды запущены
   kubectl get pods
   
4. Делаем проверочный запрос. Адрес узла можно узнать командой kubectl get nodes -o wide
   curl -X POST 192.168.106.4:30080/hello-k8s -d 'Hello, Helm!'
