#!/bin/bash

URL="http://localhost:8080/api"
HDR="Content-Type: application/json"

echo "📌 Iniciando população de dados..."

##################################################
# Criar Usuários
##################################################
echo "👤 Criando usuários..."
declare -a usuarios=(
'{"email":"alice@example.com","senha":"123","cargo":"Admin"}'
'{"email":"bruno@example.com","senha":"abc","cargo":"Técnico"}'
'{"email":"carla@example.com","senha":"123","cargo":"Estagiário"}'
'{"email":"diego@example.com","senha":"123","cargo":"Administrador"}'
'{"email":"eduarda@example.com","senha":"pass","cargo":"Cientista"}'
'{"email":"felipe@example.com","senha":"123","cargo":"Técnico"}'
'{"email":"gabriela@example.com","senha":"123","cargo":"Analista"}'
'{"email":"hugo@example.com","senha":"123","cargo":"Suporte"}'
'{"email":"isabela@example.com","senha":"123","cargo":"Estagiário"}'
'{"email":"joao@example.com","senha":"123","cargo":"Admin"}'
'{"email":"karina@example.com","senha":"123","cargo":"Cientista"}'
'{"email":"lucas@example.com","senha":"123","cargo":"Supervisor"}'
'{"email":"marina@example.com","senha":"123","cargo":"Analista"}'
'{"email":"nicolas@example.com","senha":"123","cargo":"Técnico"}'
'{"email":"olivia@example.com","senha":"123","cargo":"Gestora"}'
'{"email":"paulo@example.com","senha":"123","cargo":"Administrador"}'
'{"email":"quezia@example.com","senha":"123","cargo":"Engenheira"}'
'{"email":"renato@example.com","senha":"123","cargo":"Estagiário"}'
'{"email":"sofia@example.com","senha":"123","cargo":"Técnico"}'
'{"email":"thiago@example.com","senha":"123","cargo":"Supervisor"}'
)

for u in "${usuarios[@]}"; do
  curl -s -X POST "$URL/usuarios" -H "$HDR" -d "$u" > /dev/null
done
echo "✔ Usuários criados!"


##################################################
# Criar Almoxarifado
##################################################
echo "📦 Criando itens do almoxarifado..."

declare -a almox=(
'{"email":"Máscara N95","categoria":"EPI","data_validade":"2026-10-10","criado_por":1}'
'{"email":"Álcool em Gel","categoria":"Higiene","data_validade":"2025-05-15","criado_por":1}'
'{"email":"Luvas Cirúrgicas","categoria":"EPI","data_validade":"2027-01-21","criado_por":1}'
'{"email":"Seringas","categoria":"Materiais Médicos","criado_por":1}'
'{"email":"Bandagens","categoria":"Curativos","data_validade":"2027-09-09","criado_por":1}'
'{"email":"Esparadrapo","categoria":"Curativos","data_validade":"2027-03-15","criado_por":1}'
'{"email":"Tiras Reagentes","categoria":"Diagnóstico","data_validade":"2026-01-12","criado_por":1}'
'{"email":"Desinfetante","categoria":"Limpeza","data_validade":"2025-12-30","criado_por":1}'
'{"email":"Termômetros","categoria":"Instrumentos Médicos","criado_por":1}'
'{"email":"Gaze Estéril","categoria":"Curativos","data_validade":"2027-05-05","criado_por":1}'
'{"email":"Agulhas 30G","categoria":"Materiais Médicos","data_validade":"2027-09-01","criado_por":1}'
'{"email":"Máscara Cirúrgica","categoria":"EPI","data_validade":"2026-03-18","criado_por":1}'
'{"email":"Fita Microporosa","categoria":"Curativos","data_validade":"2027-07-22","criado_por":1}'
'{"email":"Algodão","categoria":"Higiene","criado_por":1}'
'{"email":"Kit Primeiros Socorros","categoria":"Emergência","criado_por":1}'
'{"email":"Swab Estéril","categoria":"Diagnóstico","data_validade":"2026-10-20","criado_por":1}'
'{"email":"Bisturi","categoria":"Cirúrgico","criado_por":1}'
'{"email":"Protetor Auricular","categoria":"EPI","data_validade":"2028-08-08","criado_por":1}'
'{"email":"Lençol Hospitalar","categoria":"Têxtil","criado_por":1}'
'{"email":"Água Oxigenada","categoria":"Higiene","data_validade":"2026-09-11","criado_por":1}'
)

for a in "${almox[@]}"; do
  curl -s -X POST "$URL/almoxarifado" -H "$HDR" -d "$a" > /dev/null
done
echo "✔ Almoxarifado populado!"


##################################################
# Criar Patrimônios
##################################################
echo "🏢 Criando patrimônios..."
for i in $(seq 1 20); do
  curl -s -X POST "$URL/patrimonios" \
  -H "$HDR" \
  -d "{\"email\":\"Equipamento $i\",\"identificacao_fisica\":\"PAT-$i\",\"localizacao\":\"Sala $((RANDOM%10+1))\",\"status\":\"ativo\",\"criado_por\":1}" \
  > /dev/null
done
echo "✔ Patrimônios criados!"


echo "🎯 População concluída com sucesso!"
