import MiniSearch from 'minisearch'
import recipesData from './data/recipes.json'

// ── Index ─────────────────────────────────────────────────

const miniSearch = new MiniSearch({
  fields: ['searchText'],
  storeFields: ['name', 'ingredients', 'method', 'garnish', 'imageURL', 'youtubeLink'],
  processTerm: term => term.toLowerCase(),
})

const docs = recipesData.map(r => ({
  id: r.name,
  searchText: [r.name, ...r.ingredients.map(i => i.description)].join(' '),
  name: r.name,
  ingredients: r.ingredients,
  method: r.method,
  garnish: r.garnish,
  imageURL: r.image_url,
  youtubeLink: r.youtube_link,
}))

miniSearch.addAll(docs)

// ── State ─────────────────────────────────────────────────

/** @type {{ term: string, exclude: boolean }[]} */
let chips = []

// ── Search ────────────────────────────────────────────────

function applySearch() {
  const included = chips.filter(c => !c.exclude).map(c => c.term)
  const excluded = chips.filter(c => c.exclude).map(c => c.term)

  let results = docs

  if (included.length > 0) {
    const idSets = included.map(term => {
      const hits = miniSearch.search(term, { prefix: true, fuzzy: 0.15 })
      return new Set(hits.map(h => h.id))
    })
    const intersection = idSets.reduce((a, b) => new Set([...a].filter(id => b.has(id))))
    results = docs.filter(d => intersection.has(d.id))
  }

  if (excluded.length > 0) {
    results = results.filter(d => {
      const text = d.searchText.toLowerCase()
      return !excluded.some(term => {
        const hits = miniSearch.search(term, { prefix: true, fuzzy: 0.1 })
        return hits.some(h => h.id === d.id)
      })
    })
  }

  renderGrid(results)
}

// ── Render grid ───────────────────────────────────────────

function formatAmount(ingredient) {
  if (!ingredient.amount) return ingredient.scale || ''
  const amount = ingredient.amount % 1 === 0
    ? ingredient.amount.toString()
    : ingredient.amount.toFixed(1)
  return ingredient.scale ? `${amount} ${ingredient.scale}` : amount
}

function ingredientPreview(ingredients) {
  return ingredients
    .slice(0, 4)
    .map(i => i.description)
    .join(', ') + (ingredients.length > 4 ? ` +${ingredients.length - 4} more` : '')
}

function renderGrid(results) {
  const grid = document.getElementById('grid')
  const empty = document.getElementById('empty')
  const count = document.getElementById('result-count')

  if (results.length === 0) {
    grid.innerHTML = ''
    empty.hidden = false
    count.textContent = ''
    return
  }

  empty.hidden = true
  count.textContent = `${results.length} cocktail${results.length === 1 ? '' : 's'}`

  grid.innerHTML = results.map(r => `
    <article class="card" data-id="${encodeURIComponent(r.id)}" tabindex="0" role="button" aria-label="${r.name}">
      ${r.imageURL
        ? `<img class="card-img" src="${r.imageURL}" alt="${r.name}" loading="lazy" onerror="this.replaceWith(makePlaceholder())">`
        : `<div class="card-img-placeholder">🍹</div>`}
      <div class="card-body">
        <p class="card-name">${r.name}</p>
        <p class="card-ingredients">${ingredientPreview(r.ingredients)}</p>
      </div>
    </article>
  `).join('')

  grid.querySelectorAll('.card').forEach(card => {
    const open = () => openModal(decodeURIComponent(card.dataset.id))
    card.addEventListener('click', open)
    card.addEventListener('keydown', e => { if (e.key === 'Enter' || e.key === ' ') open() })
  })
}

// ── Modal ─────────────────────────────────────────────────

function openModal(id) {
  const recipe = docs.find(d => d.id === id)
  if (!recipe) return

  const body = document.getElementById('modal-body')
  const garnishHtml = recipe.garnish && recipe.garnish !== 'N/A'
    ? `<p class="modal-garnish">🌿 Garnish: <span>${recipe.garnish}</span></p>`
    : ''

  const ytHtml = recipe.youtubeLink
    ? `<a class="btn-link btn-yt" href="${recipe.youtubeLink}" target="_blank" rel="noopener">▶ Watch on YouTube</a>`
    : ''

  body.innerHTML = `
    ${recipe.imageURL
      ? `<img class="modal-img" src="${recipe.imageURL}" alt="${recipe.name}">`
      : `<div class="modal-img-placeholder">🍹</div>`}
    <h2 class="modal-name">${recipe.name}</h2>

    <p class="modal-section-title">Ingredients</p>
    <ul class="modal-ingredients">
      ${recipe.ingredients.map(i => `
        <li class="modal-ingredient">
          <span class="ingredient-amount">${formatAmount(i)}</span>
          <span class="ingredient-name">${i.description}</span>
        </li>
      `).join('')}
    </ul>

    <p class="modal-section-title">Method</p>
    <p class="modal-text">${recipe.method}</p>

    ${garnishHtml}

    ${ytHtml ? `<div class="modal-links">${ytHtml}</div>` : ''}
  `

  document.getElementById('modal').showModal()
}

// ── Chips ─────────────────────────────────────────────────

function renderChips() {
  const container = document.getElementById('chips')
  container.innerHTML = chips.map((c, i) => `
    <span class="chip ${c.exclude ? 'exclude' : 'include'}">
      ${c.exclude ? '−' : '+'} ${c.term}
      <button class="chip-remove" data-index="${i}" aria-label="Remove ${c.term}">✕</button>
    </span>
  `).join('')

  container.querySelectorAll('.chip-remove').forEach(btn => {
    btn.addEventListener('click', () => {
      chips.splice(Number(btn.dataset.index), 1)
      renderChips()
      applySearch()
    })
  })
}

function addChip(raw) {
  const term = raw.trim().replace(/^,+|,+$/g, '').trim()
  if (!term) return
  const exclude = term.startsWith('-')
  const cleanTerm = exclude ? term.slice(1).trim() : term
  if (!cleanTerm) return
  if (chips.some(c => c.term.toLowerCase() === cleanTerm.toLowerCase() && c.exclude === exclude)) return
  chips.push({ term: cleanTerm, exclude })
  renderChips()
  applySearch()
}

// ── Event wiring ──────────────────────────────────────────

document.getElementById('search-input').addEventListener('keydown', e => {
  const input = e.currentTarget
  if (e.key === 'Enter' || e.key === ',') {
    e.preventDefault()
    addChip(input.value)
    input.value = ''
  } else if (e.key === 'Backspace' && input.value === '' && chips.length > 0) {
    chips.pop()
    renderChips()
    applySearch()
  }
})

document.getElementById('chips-input').addEventListener('click', () => {
  document.getElementById('search-input').focus()
})

document.getElementById('modal-close').addEventListener('click', () => {
  document.getElementById('modal').close()
})

document.getElementById('modal').addEventListener('click', e => {
  if (e.target === e.currentTarget) e.currentTarget.close()
})

// ── Init ──────────────────────────────────────────────────

renderGrid(docs)
document.getElementById('result-count').textContent = `${docs.length} cocktails`
