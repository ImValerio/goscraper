import { useEffect, useState, useContext, createContext } from "react";
import Result from "./components/Result";
import Setups from "./components/Setups";
import Modal from "./components/Modal";
import "./App.css";

const AppContext = createContext();

export const useAppContext = () => useContext(AppContext);

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState([]);
  const [setups, setSetups] = useState([]);
  const [urls, setUrls] = useState("");
  const [pattern, setPattern] = useState("");

  const [modalData, setModalData] = useState(null);
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    loadSetups();
  }, []);

  const loadSetups = () => {
    const setups = window.localStorage.getItem("setups");
    if (setups && setups !== "") {
      setSetups(JSON.parse(setups));
    }
  };

  const scrapeUrls = async () => {
    setIsLoading(true);

    const urlsToScrape = urls.split("\n");
    const tags = pattern.split("\n").map((el) => el.split("->"));

    const serverHost = process.env.SERVER_HOST || "localhost";
    const serverPort = process.env.SERVER_PORT || "5000";
    const res = await fetch(`http://${serverHost}:${serverPort}/scrape`, {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({ urls: urlsToScrape, tags }),
    });

    const data = await res.json();

    setResult(data);

    setIsLoading(false);
  };

  const saveSetup = () => {
    if (urls.length <= 0 && pattern === "") return;

    const setup = { urls, pattern };

    window.localStorage.setItem("setups", JSON.stringify([...setups, setup]));
    setSetups([...setups, setup]);
  };

  const removeSetup = (urls, pattern) => {
    const newSetups = setups.filter(
      (setup) => setup.urls !== urls && setup.pattern !== pattern
    );

    setSetups(newSetups);

    window.localStorage.setItem("setups", JSON.stringify(newSetups));
  };

  const downloadCsv = () => {
    const rows = [["url", "data"]];

    rows.push(...result.map((el) => [el.url, el.res.join("|")]));

    let csvContent =
      "data:text/csv;charset=utf-8," + rows.map((e) => e.join(",")).join("\n");

    var encodedUri = encodeURI(csvContent);
    var link = document.createElement("a");
    link.setAttribute("href", encodedUri);
    link.setAttribute("download", `goscraper_${Date.now()}.csv`);
    document.body.appendChild(link); // Required for FF

    link.click();
  };

  return (
    <AppContext.Provider
      value={{
        urls,
        setUrls,
        pattern,
        setPattern,
        removeSetup,
        showModal,
        setShowModal,
        setModalData,
      }}
    >
      <div className="w-full h-full flex justify-center items-center flex-col">
        <h1 className="text-5xl font-bold my-3">
          <span className="text-cyan-500">Go</span>Scraper
        </h1>
        <div className="flex flex-col w-screen max-w-screen-lg px-5 my-3">
          {setups && <Setups items={setups} />}
          <textarea
            className="text-xl my-1 border border-1 p-1"
            placeholder="https://en.wikipedia.org/wiki/Computer_programming"
            value={urls}
            onChange={(e) => setUrls(e.target.value)}
          ></textarea>
          <textarea
            className="text-xl my-1 border border-1 p-1"
            placeholder="div->h1->a"
            value={pattern}
            onChange={(e) => setPattern(e.target.value)}
          ></textarea>
          <div className="flex w-100">
            <button
              className={
                result.length > 0
                  ? "text-2xl font-bold px-2 py-1 bg-green-500 hover:bg-green-700 text-white transition-all"
                  : "text-2xl font-bold px-2 py-1 bg-gray-500 hover:bg-gray-700 text-white transition-all cursor-not-allowed"
              }
              onClick={() => downloadCsv()}
              disabled={result.length <= 0}
              title="Download CSV"
            >
              <img src="download.svg" className="svg-white" />
            </button>
            <button
              className="text-2xl font-bold px-2 py-1 bg-cyan-500 hover:bg-cyan-700 text-white transition-all grow mx-2"
              onClick={() => scrapeUrls()}
            >
              SCRAPE
            </button>
            <button
              className="text-2xl font-bold px-2 py-1 bg-green-500 hover:bg-green-700 text-white transition-all"
              onClick={() => saveSetup()}
            >
              <img src="save.svg" className="svg-white" />
            </button>
          </div>
        </div>
        <div className="flex">
          {result.map((el) => (
            <Result element={el} />
          ))}
        </div>
        {showModal && <Modal data={modalData} setShowModal={setShowModal} />}
      </div>
    </AppContext.Provider>
  );
}

export default App;
