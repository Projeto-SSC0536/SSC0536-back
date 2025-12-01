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
'{"nome":"Máscara N95","categoria":"EPI","data_validade":"2026-10-10T00:00:00Z","criado_por":1}'
'{"nome":"Álcool em Gel","categoria":"Higiene","data_validade":"2025-05-15T00:00:00Z","criado_por":1}'
'{"nome":"Luvas Cirúrgicas","categoria":"EPI","data_validade":"2027-01-21T00:00:00Z","criado_por":1}'
'{"nome":"Seringas","categoria":"Materiais Médicos","criado_por":1}'
'{"nome":"Bandagens","categoria":"Curativos","data_validade":"2027-09-09T00:00:00Z","criado_por":1}'
'{"nome":"Esparadrapo","categoria":"Curativos","data_validade":"2027-03-15T00:00:00Z","criado_por":1}'
'{"nome":"Tiras Reagentes","categoria":"Diagnóstico","data_validade":"2026-01-12T00:00:00Z","criado_por":1}'
'{"nome":"Desinfetante","categoria":"Limpeza","data_validade":"2025-12-30T00:00:00Z","criado_por":1}'
'{"nome":"Termômetros","categoria":"Instrumentos Médicos","criado_por":1}'
'{"nome":"Gaze Estéril","categoria":"Curativos","data_validade":"2027-05-05T00:00:00Z","criado_por":1}'
'{"nome":"Agulhas 30G","categoria":"Materiais Médicos","data_validade":"2027-09-01T00:00:00Z","criado_por":1}'
'{"nome":"Máscara Cirúrgica","categoria":"EPI","data_validade":"2026-03-18T00:00:00Z","criado_por":1}'
'{"nome":"Fita Microporosa","categoria":"Curativos","data_validade":"2027-07-22T00:00:00Z","criado_por":1}'
'{"nome":"Algodão","categoria":"Higiene","criado_por":1}'
'{"nome":"Kit Primeiros Socorros","categoria":"Emergência","criado_por":1}'
'{"nome":"Swab Estéril","categoria":"Diagnóstico","data_validade":"2026-10-20T00:00:00Z","criado_por":1}'
'{"nome":"Bisturi","categoria":"Cirúrgico","criado_por":1}'
'{"nome":"Protetor Auricular","categoria":"EPI","data_validade":"2028-08-08T00:00:00Z","criado_por":1}'
'{"nome":"Lençol Hospitalar","categoria":"Têxtil","criado_por":1}'
'{"nome":"Água Oxigenada","categoria":"Higiene","data_validade":"2026-09-11T00:00:00Z","criado_por":1}'
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
  -d "{\"nome\":\"Equipamento $i\",\"identificacao_fisica\":\"PAT-$i\",\"localizacao\":\"Sala $((RANDOM%10+1))\",\"status\":\"ativo\",\"criado_por\":1}" \
  > /dev/null
done
echo "✔ Patrimônios criados!"


echo "🎯 População concluída com sucesso!"
