const h = (x) => document.createElement(x);
const t = (x) => document.createTextNode(x);
function a(link, text) {
    const ele = h('a');
    ele.href = link;
    ele.text = text;
    ele.target = '_blank';
    return ele;
}
const DATA_URL = './data.json';
const sort = {
    current: { key: null, dir: null },
    previous: { key: null, dir: null }
};
const tablePointers = {
    thead: null,
    theadCells: {},
    tbody: null,
    tbodyRows: [],
};
const arrows = {
    updown: '⇅',
    up: '▲',
    down: '▼',
};
const columns = {
    '': '#',
    language: 'language',
    year: 'year',
    stars: 'stars',
    size_mb: 'size (mb)',
    ram_kb: 'ram (kb)',
    ips: 'iter/sec',
};
let data;
const response = fetch(DATA_URL);
const json = response.then(x => x.json());
json.then(x => {
    data = x;
    renderTable();
});
function handleSort(newSortKey) {
    updateSort(newSortKey);
    sortData();
    updateTable();
}
function updateSort(newSortKey) {
    // calculate new sort
    let newSortDir = null;
    if (sort.current.key !== newSortKey) {
        newSortDir = 'asc';
    }
    else {
        if (sort.current.dir === 'asc') {
            newSortDir = 'desc';
        }
        else {
            newSortDir = 'asc';
        }
    }
    // update sort
    sort.previous.key = sort.current.key;
    sort.previous.dir = sort.current.dir;
    sort.current.key = newSortKey;
    sort.current.dir = newSortDir;
}
function sortData() {
    const sortKey = sort.current.key;
    const sortDir = sort.current.dir;
    if (!sortKey)
        return;
    tablePointers.tbodyRows.sort((a, b) => {
        const first = a.data[sortKey];
        const second = b.data[sortKey];
        if (first === second) {
            return 0;
        }
        else if (!first) {
            return 1;
        }
        else if (!second) {
            return -1;
        }
        else if (first < second) {
            return (sortDir === 'asc') ? -1 : 1;
        }
        else {
            return (sortDir === 'asc') ? 1 : -1;
        }
    });
}
function updateTable() {
    const previousSortKey = sort.previous.key;
    const sortKey = sort.current.key;
    const sortDir = sort.current.dir;
    const theadCells = tablePointers.theadCells;
    // =========================
    // clear previous header
    if (previousSortKey && previousSortKey !== sortKey) {
        const previousDisplay = columns[previousSortKey];
        theadCells[previousSortKey].data = `${previousDisplay} ${arrows.updown}`;
    }
    // set current header
    const currentDisplay = columns[sortKey];
    const modifier = (sortDir === 'asc') ? arrows.up : arrows.down;
    theadCells[sortKey].data = `${currentDisplay} ${modifier}`;
    // update index column
    for (const [i, row] of tablePointers.tbodyRows.entries()) {
        row.idxCellText.data = (i + 1).toString();
    }
    // update table
    tablePointers.tbody.replaceChildren(...tablePointers.tbodyRows.map(x => x.ele));
}
function renderTable() {
    if (!data || !data.length)
        return;
    const table = h('table');
    // head
    const thead = h('thead');
    tablePointers.thead = thead;
    const theadRow = h('tr');
    for (const [name, display] of Object.entries(columns)) {
        console.log(name, display);
        const cell = h('th');
        const text = name ? t(`${display} ⇅`) : t(display);
        if (name) {
            cell.addEventListener('click', () => handleSort(name));
        }
        tablePointers.theadCells[name] = text;
        cell.append(text);
        theadRow.append(cell);
    }
    thead.append(theadRow);
    table.append(thead);
    // body
    const tbody = h('tbody');
    tablePointers.tbody = tbody;
    for (const [i, x] of data.entries()) {
        const trow = h('tr');
        const cellContents = [
            t((i + 1).toString()),
            t(x.language),
            x.wiki_url ? a(x.wiki_url, x.year.toString()) : t(x.year.toString()),
            x.github_url ? a(x.github_url, x.stars.toLocaleString()) : t(""),
            t(Math.round(x.size_mb).toLocaleString()),
            t(Math.round(x.ram_kb).toLocaleString()),
            t(x.ips.toLocaleString())
        ];
        for (const inner of cellContents) {
            const cell = h('td');
            cell.append(inner);
            trow.append(cell);
        }
        tablePointers.tbodyRows.push({ ele: trow, data: x, idxCellText: cellContents[0] });
        tbody.append(trow);
    }
    table.append(tbody);
    // append
    document.getElementById('table-container').replaceChildren(table);
}
export {};
//# sourceMappingURL=script.js.map