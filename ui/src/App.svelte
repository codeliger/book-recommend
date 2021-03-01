<script>

	import {onMount} from "svelte"

	let serverURL = "http://127.0.0.1:8080/v1/"

	let allAuthors = []
	let allGenres = []
	let selectedAuthors = []
	let selectedGenres = []
	let minPageCountRange = 0
	let maxPageCountRange = 0
	let minPageCount = 0
	let maxPageCount = 0

	let books = []

	async function populateFilters(){
		console.log("mounted")
		let response = await fetch(serverURL + "filters")
		if (response.ok && response.status == 200) {
			let filters = await response.json()
			console.log(filters)
			allAuthors = filters.authors
			allGenres = filters.genres
			minPageCountRange = filters.min_page_count
			maxPageCountRange = filters.max_page_count
			minPageCount = minPageCountRange
			maxPageCount = maxPageCountRange
		}else{
			console.log("failed", response.statusText, await response.text())
		}
	}

	function toggleGenre(){

		let genre = this.getAttribute("data-id")
		console.log('genre from data-id', genre)

		if (selectedGenres.includes(genre)){
			selectedGenres.splice(selectedGenres.indexOf(genre), 1)
			console.log("removing", genre)
		}else{
			selectedGenres.push(genre)
			console.log("adding", genre)
		}

		selectedGenres = selectedGenres
	}

	function toggleAuthor(){
		let author = this.getAttribute("data-id")
		console.log('author from data-id', author)

		if (selectedAuthors.includes(author)){
			selectedAuthors.splice(selectedAuthors.indexOf(author), 1)
			console.log("removing", author)
		}else{
			selectedAuthors.push(author)
			console.log("adding", author)
		}

		selectedAuthors = selectedAuthors
	}
	
	onMount(async () => {
		await populateFilters()
	})

	async function search() {

		let data = {"min_page_count": minPageCount, "max_page_count": maxPageCount, "genres": selectedGenres, "authors": selectedAuthors}

		let resp = await fetch(serverURL + `ranked_books?min_page_count=${minPageCount}&max_page_count=${maxPageCount}&genres=${selectedGenres.join(":")}&authors=${selectedAuthors.join(":")}`)

		if (resp.ok && resp.status == 200){
			books = await resp.json()
		}
	}
</script> 

<main>
	<div class='wrapper'>
		<div class='categories'>
			<div class='authors'>
				<h2>Authors</h2>
				{#each allAuthors as author, i}
					<button data-id="{i}" on:click={toggleAuthor} class:selected={selectedAuthors.includes(String(i))}>{author}</button>
				{/each}
			</div>
			<div class='genres'>
				<h2>Genres</h2>
				{#each allGenres as genre, i}
					<button data-id="{i}" on:click={toggleGenre} class:selected={selectedGenres.includes(String(i))}>{genre}</button>
				{/each}
			</div>
			<div class="min_pages">
				<h2>Minimum Pages</h2>
				<input type=range min={minPageCountRange} max={maxPageCountRange} bind:value={minPageCount}>
				{minPageCount}
			</div>
			<div class="max_pages">
				<h2>Maximum Pages</h2>
				<input type=range min={minPageCountRange} max={maxPageCountRange} bind:value={maxPageCount}>
				{maxPageCount}
			</div>
			<button id="btnSearch" on:click={search}>Search</button>
		</div>
		<div class='books'>
				{#each books as book}
				<div class="book">
					<h4>{book.title}</h4>
					<p>Authors: {book.authors.join(",")}</p>
					<p>Genres: {book.genres.join(",")}</p>
					<p>Pages: {book.pages}</p>
					<p>Publication Year: {book.publication_year}</p>
					<p>Rating: {book.rating}</p>
				</div>
				{/each}
		</div>
	</div>
</main>

<style>
	main {
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}

	.selected{
		background-color: blue;
		color: white;
	}

	.wrapper{
		display:flex;
	}

	.categories {
		width:50vw;
	}

	.books {
		width:50vw;
		display:flex;
		flex-direction: row;
		flex-wrap: wrap;
	}

	.book{
		margin:15px;
		padding:5px;
	}

	.book h4{
		max-width:200px;
	}


</style>