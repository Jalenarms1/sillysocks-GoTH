const listCont = document.getElementById("cart-items-list")

		const checkAllBtn = document.getElementById("check-all")
		const delItemsBtn = document.getElementById("delete-items")
		const writeToDel = document.getElementById("del-ids")
		if (checkAllBtn) {

			checkAllBtn.addEventListener("click", (e) => {
				if (e.target.checked) {
					delItemsBtn.classList.remove("hidden")
				} else {
					delItemsBtn.classList.add("hidden")
					writeToDel.value = ""

				}
				const chbxs = document.querySelectorAll("input[type='checkbox']")
				chbxs.forEach(c => {
					if (c.id != "check-all") {
						c.checked = e.target.checked
						const currIds = writeToDel.value.split(",")
						if (e.target.checked) {
							if (!currIds.includes(c.id)) {
								currIds.push(c.id)
								writeToDel.value = currIds.join(",")

							}

						} 

					}
				})
			})
		}
		listCont.addEventListener("click", (e) => {
			if (e.target.classList.contains("checkbox")) {
				const cartItemId = e.target.id
				const currIds = writeToDel.value.split(",")
				if (e.target.checked) {
					if (!currIds.includes(cartItemId)) {
						currIds.push(cartItemId)
						writeToDel.value = currIds.join(",")

					}
				} else {
					const newIds = currIds.filter(id => id != cartItemId)

					writeToDel.value = newIds.join(",")
				}
				if (writeToDel.value == "") {
					delItemsBtn.classList.add("hidden")
				} else {
					delItemsBtn.classList.remove("hidden")
				}
			}
		})