package services

var cnaeToCategory = map[int]string{
	9999999: "Outros",
	4711302: "Supermercados",
	4721102: "Padarias",
	4771701: "Farmácias",
	4722901: "Açougues",
	5611201: "Restaurantes",
	5611203: "Lanchonetes",
	4781400: "Lojas de roupas",
	4761003: "Papelarias",
	4761001: "Livrarias",
	4774100: "Óticas",
	4783102: "Relojoarias",
	4763601: "Brinquedos",
	4754701: "Lojas de móveis",
	4751201: "Lojas de informática",
	4752100: "Lojas de telefonia",
	4731800: "Postos de combustíveis",
	9602501: "Cabeleireiros",
	4763602: "Artigos esportivos",
	4782201: "Lojas de calçados",
	4783101: "Joalherias",
	4729602: "Lojas de conveniência",
}

func MapCNAEToCategory(cnae int) string {
	if name, ok := cnaeToCategory[cnae]; ok {
		return name
	}
	return "Outros"
}
