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

      // Template for the Division's popup
      const template = {
        title: "USACE Civil Works Division",
        content: "{DIV_SYM} : {DIVISION}",
      };
      // Layer that pulls in the geoJSON data for the USACE Divisions
      const geoLayer = new GeoJSONLayer({
        url: "https://opendata.arcgis.com/datasets/4cdab8820d7a4f58aaa003be63f059ac_0.geojson",
        copyright: "USACE Civil Works",
        popupTemplate: template,
        title: "USACE Divisions",
        listMode: false
      });

      // Creates the map component
      const map = new ArcGISMap({
        basemap: "topo-vector",  // initial map styling
        layers: [geoLayer]       // array of layers that sits on top of the basemap
      });

      // Creates the MapView - Necessary for rendering the Map Object above
      const view = new MapView({
        map: map,
        container: mapDiv.current,
        // sets the initial positioning and zoom of the map
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

      //Creates the LayerList which is used to toggle visiblity of available map layers
      view.when(function () {
        var layerList = new LayerList({
          view: view,
        });
        // adds the Layer toggle to the top-right of the screen
        view.ui.add(layerList, "top-right");
      });

      //BasemapToggle object is used to toggle the "style" of the basemap 
      var toggle = new BasemapToggle({
        view: view, // view that provides access to the map's 'topo-vector' basemap
        nextBasemap: "hybrid" // allows for toggling to the 'hybrid' basemap
      });
      // adds the toggle base map button display
      view.ui.add(toggle, "top-right");
    }
  }, []);

  return (<div className="mapDiv" ref={mapDiv}></div>

  );
}

export default App;