const h = (x) => document.createElement(x);
const t = (x) => document.createTextNode(x);
function a(link, text) {
    const ele = h('a');
    ele.href = link;
    ele.text = text;
    ele.target = '_blank';
    return ele;
}
const global = {
    fileName: './data.json',
    data: null,
    sort: {
        current: { key: null, dir: null },
        previous: { key: null, dir: null }
    },
    columns: {
        language: { display: 'language', createNode: (x) => t(x.language) },
        year: { display: 'year', createNode: (x) => x.wiki_url ? a(x.wiki_url, x.year.toString()) : t(x.year.toString()) },
        stars: { display: 'stars', createNode: (x) => x.github_url ? a(x.github_url, x.stars.toLocaleString()) : t("") },
        size_mb: { display: 'size (mb)', createNode: (x) => t(Math.round(x.size_mb).toLocaleString()) },
        ram_kb: { display: 'ram (kb)', createNode: (x) => t(Math.round(x.ram_kb).toLocaleString()) },
        ips: { display: 'ips', createNode: (x) => t(x.ips.toLocaleString()) }
    },
    tablePointers: {
        thead: null,
        theadCells: {},
        tbody: null,
        tbodyRows: [],
    },
    arrows: {
        updown: '⇅',
        up: '▲',
        down: '▼',
    },
    functions: {
        getData: null,
        handleSort: null,
        updateSort: null,
        sortData: null,
        updateTable: null,
        renderTable: null,
    }
};
window.global = global;
// initialization
(function fetchData() {
    //
    Object.assign(global.functions, { getData, handleSort, updateSort, sortData, updateTable, renderTable });
    // externals
    const setData = (x) => global.data = x;
    // ========================================
    const response = fetch('./data.json');
    const json = response.then(x => x.json());
    json.then(x => {
        setData(x);
        renderTable();
    });
    // functions
    function getData() {
        // externals
        const fileName = global.fileName;
        const setData = (x) => global.data = x;
        // ========================================
        const response = fetch(fileName);
        const json = response.then(x => x.json());
        return json.then(x => {
            setData(x);
            renderTable();
        });
    }
    function handleSort(newSortKey) {
        // externals
        const updateSort = global.functions.updateSort;
        const sortData = global.functions.sortData;
        const updateTable = global.functions.updateTable;
        // ==============================
        updateSort(newSortKey);
        sortData();
        updateTable();
    }
    function updateSort(newSortKey) {
        // externals
        const oldSortKey = global.sort.current.key;
        const oldSortDir = global.sort.current.dir;
        const setSortKey = (x) => global.sort.current.key = x;
        const setSortDir = (x) => global.sort.current.dir = x;
        const setPreviousSort = ({ key, dir }) => {
            global.sort.previous.key = key;
            global.sort.previous.dir = dir;
        };
        // ==============================
        // calculate new sort
        let newSortDir = null;
        if (oldSortKey !== newSortKey) {
            newSortDir = 'asc';
        }
        else {
            if (oldSortDir === 'asc') {
                newSortDir = 'desc';
            }
            else {
                newSortDir = 'asc';
            }
        }
        // update sort
        setPreviousSort({ key: oldSortKey, dir: oldSortDir });
        setSortKey(newSortKey);
        setSortDir(newSortDir);
    }
    function sortData() {
        // externals
        const sortKey = global.sort.current.key;
        const sortDir = global.sort.current.dir;
        const sortable = global.tablePointers.tbodyRows;
        // ========================================
        if (!sortKey)
            return;
        sortable.sort((a, b) => {
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
        // externals
        const previousSortKey = global.sort.previous.key;
        const sortKey = global.sort.current.key;
        const sortDir = global.sort.current.dir;
        const tbody = global.tablePointers.tbody;
        const rows = global.tablePointers.tbodyRows;
        const theadCells = global.tablePointers.theadCells;
        const updown = global.arrows.updown;
        const up = global.arrows.up;
        const down = global.arrows.down;
        const columns = global.columns;
        // =========================
        if (previousSortKey && previousSortKey !== sortKey) {
            const previousDisplay = columns[previousSortKey].display;
            theadCells[previousSortKey].data = `${previousDisplay} ${updown}`;
        }
        const currentDisplay = columns[sortKey].display;
        const modifier = (sortDir === 'asc') ? up : down;
        theadCells[sortKey].data = `${currentDisplay} ${modifier}`;
        tbody.replaceChildren(...rows.map(x => x.ele));
    }
    function renderTable() {
        // externals
        const localData = global.data;
        const setThead = (x) => global.tablePointers.thead = x;
        const setTbody = (x) => global.tablePointers.tbody = x;
        const localTheadMap = global.tablePointers.theadCells;
        const localTbodyRows = global.tablePointers.tbodyRows;
        const columns = Object.entries(global.columns);
        // ========================================
        const h = (x) => document.createElement(x);
        const t = (x) => document.createTextNode(x);
        if (!localData || !localData.length)
            return;
        const table = h('table');
        // head
        const thead = h('thead');
        setThead(thead);
        const theadRow = h('tr');
        for (const [key, { display }] of columns) {
            const cell = h('th');
            const text = t(`${display} ⇅`);
            cell.addEventListener('click', () => handleSort(key));
            localTheadMap[key] = text;
            cell.append(text);
            theadRow.append(cell);
        }
        thead.append(theadRow);
        table.append(thead);
        // body
        const tbody = h('tbody');
        setTbody(tbody);
        for (const dataRow of localData) {
            const trow = h('tr');
            for (const [_, { createNode }] of columns) {
                const cell = h('td');
                const node = createNode(dataRow);
                cell.append(node);
                trow.append(cell);
            }
            localTbodyRows.push({ ele: trow, data: dataRow });
            tbody.append(trow);
        }
        table.append(tbody);
        // append
        document.getElementById('table-container').replaceChildren(table);
    }
})();
export {};
//# sourceMappingURL=script.js.map