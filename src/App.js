import React, { useRef, useEffect } from "react";
//import FeatureLayer from "@arcgis/core/layers/FeatureLayer";
import ArcGISMap from "@arcgis/core/Map";
//import DictionaryRenderer from "@arcgis/core/renderers/DictionaryRenderer";
import MapView from "@arcgis/core/views/MapView";
import GeoJSONLayer from "@arcgis/core/layers/GeoJSONLayer";
import LayerList from "@arcgis/core/widgets/LayerList";
import BasemapToggle from "@arcgis/core/widgets/BasemapToggle";
import esriConfig from '@arcgis/core/config.js';

import "./App.css";

function App() {

  // Required: Set this property to insure assets resolve correctly.
  esriConfig.assetsPath = './assets';

  const mapDiv = useRef(null);

  useEffect(() => {
    if (mapDiv.current) {
      /**
       * Initialize application
       */

      const template = {
        title: "USACE Civil Works Division",
        content: "{DIVISION} {DIV_SYM}",
      };

      const geoLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/4cdab8820d7a4f58aaa003be63f059ac_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: template,
        title: "USACE Divisions",
        listMode: false
      });

      const map = new ArcGISMap({
        basemap: "topo-vector",
        layers: [geoLayer]
      });

      const view = new MapView({
        map: map,
        container: mapDiv.current,
        extent: {
          spatialReference: {
            wkid: 102100,
          },
          xmax: -47458864.13281338,
          xmin: -54033671.557789356,
          ymax: 6678183.73641666,
          ymin: 1937864.9902844299,
        },
        zoom: 4
      });

      view.when(function () {
        var layerList = new LayerList({
          view: view,

        });

        view.ui.add(layerList, "top-right");
      });
      var toggle = new BasemapToggle({
        // 2 - Set properties
        view: view, // view that provides access to the map's 'topo-vector' basemap
        nextBasemap: "hybrid" // allows for toggling to the 'hybrid' basemap
      });
      view.ui.add(toggle, "top-right");
    }
  }, []);

  return (<div className="mapDiv" ref={mapDiv}></div>

  );
}

export default App;