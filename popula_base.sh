#!/bin/bash

URL="http://localhost:8080/api"
HDR="Content-Type: application/json"

echo "📌 Iniciando população de dados..."

##################################################
# Criar Usuários
##################################################
echo "👤 Criando usuários..."
declare -a usuarios=(
'{"nome":"Alice","senha":"123","cargo":"Admin"}'
'{"nome":"Bruno","senha":"abc","cargo":"Técnico"}'
'{"nome":"Carla","senha":"123","cargo":"Estagiário"}'
'{"nome":"Diego","senha":"123","cargo":"Administrador"}'
'{"nome":"Eduarda","senha":"pass","cargo":"Cientista"}'
'{"nome":"Felipe","senha":"123","cargo":"Técnico"}'
'{"nome":"Gabriela","senha":"123","cargo":"Analista"}'
'{"nome":"Hugo","senha":"123","cargo":"Suporte"}'
'{"nome":"Isabela","senha":"123","cargo":"Estagiário"}'
'{"nome":"João","senha":"123","cargo":"Admin"}'
'{"nome":"Karina","senha":"123","cargo":"Cientista"}'
'{"nome":"Lucas","senha":"123","cargo":"Supervisor"}'
'{"nome":"Marina","senha":"123","cargo":"Analista"}'
'{"nome":"Nicolas","senha":"123","cargo":"Técnico"}'
'{"nome":"Olivia","senha":"123","cargo":"Gestora"}'
'{"nome":"Paulo","senha":"123","cargo":"Administrador"}'
'{"nome":"Quezia","senha":"123","cargo":"Engenheira"}'
'{"nome":"Renato","senha":"123","cargo":"Estagiário"}'
'{"nome":"Sofia","senha":"123","cargo":"Técnico"}'
'{"nome":"Thiago","senha":"123","cargo":"Supervisor"}'
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
'{"nome":"Máscara N95","categoria":"EPI","data_validade":"2026-10-10","criado_por":1}'
'{"nome":"Álcool em Gel","categoria":"Higiene","data_validade":"2025-05-15","criado_por":1}'
'{"nome":"Luvas Cirúrgicas","categoria":"EPI","data_validade":"2027-01-21","criado_por":1}'
'{"nome":"Seringas","categoria":"Materiais Médicos","criado_por":1}'
'{"nome":"Bandagens","categoria":"Curativos","data_validade":"2027-09-09","criado_por":1}'
'{"nome":"Esparadrapo","categoria":"Curativos","data_validade":"2027-03-15","criado_por":1}'
'{"nome":"Tiras Reagentes","categoria":"Diagnóstico","data_validade":"2026-01-12","criado_por":1}'
'{"nome":"Desinfetante","categoria":"Limpeza","data_validade":"2025-12-30","criado_por":1}'
'{"nome":"Termômetros","categoria":"Instrumentos Médicos","criado_por":1}'
'{"nome":"Gaze Estéril","categoria":"Curativos","data_validade":"2027-05-05","criado_por":1}'
'{"nome":"Agulhas 30G","categoria":"Materiais Médicos","data_validade":"2027-09-01","criado_por":1}'
'{"nome":"Máscara Cirúrgica","categoria":"EPI","data_validade":"2026-03-18","criado_por":1}'
'{"nome":"Fita Microporosa","categoria":"Curativos","data_validade":"2027-07-22","criado_por":1}'
'{"nome":"Algodão","categoria":"Higiene","criado_por":1}'
'{"nome":"Kit Primeiros Socorros","categoria":"Emergência","criado_por":1}'
'{"nome":"Swab Estéril","categoria":"Diagnóstico","data_validade":"2026-10-20","criado_por":1}'
'{"nome":"Bisturi","categoria":"Cirúrgico","criado_por":1}'
'{"nome":"Protetor Auricular","categoria":"EPI","data_validade":"2028-08-08","criado_por":1}'
'{"nome":"Lençol Hospitalar","categoria":"Têxtil","criado_por":1}'
'{"nome":"Água Oxigenada","categoria":"Higiene","data_validade":"2026-09-11","criado_por":1}'
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

