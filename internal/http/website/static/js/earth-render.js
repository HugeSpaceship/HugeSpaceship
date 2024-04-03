import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import {CSS3DObject, CSS3DRenderer} from "three/addons/renderers/CSS3DRenderer";
import {TrackballControls} from "three/addons/controls/TrackballControls";
import {func} from "three/addons/nodes/code/FunctionNode";
const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera( 40, 1, 1, 10000 );

const renderer = new CSS3DRenderer();
renderer.setSize( window.innerWidth, window.innerHeight );
document.getElementById( 'container' ).appendChild( renderer.domElement );

document.body.appendChild( renderer.domElement );

camera.position.z = 3000;
const loader = new GLTFLoader();

loader.load( '/static/themes/shuttle/mesh/earth.gltf', function ( gltf ) {
    scene.add( gltf.scene );

}, undefined, function ( error ) {

    console.error( error );

} );

const vector = new THREE.Vector3();

for (let i = 0; i < levels.length; i += 3) {

    const element = document.createElement( 'div' );
    element.className = 'element';
    element.style.backgroundColor = 'rgba(0,127,127,' + ( Math.random() * 0.5 + 0.25 ) + ')';

    const name = document.createElement( 'div' );
    name.className = 'levelHover';
    name.textContent = levels[i];
    element.appendChild( name );

    const objectCSS = new CSS3DObject( element );
    const phi = Math.acos( - 1 + ( 2 * i ) / levels.length );
    const theta = Math.sqrt( levels.length * Math.PI ) * phi;
    // objectCSS.position.setFromSphericalCoords( 10, levels[i+1]/10, levels[i+2]/10);
    objectCSS.position.setFromSphericalCoords( 800, phi, theta);
    scene.add( objectCSS );

    const object = new THREE.Object3D();

    vector.copy( object.position ).multiplyScalar( 2 );

    objectCSS.lookAt( vector );

    // targets.sphere.push( object );
}

const light = new THREE.AmbientLight( 0xE0E0E0 ); // soft white light
scene.add( light );

const controls = new TrackballControls( camera, renderer.domElement );
controls.minDistance = 500;
controls.maxDistance = 6000;
controls.addEventListener( 'change', render );

function animate() {
    requestAnimationFrame( animate );
    // cube.rotation.x += 0.01;
    // cube.rotation.y += 0.01;
    // scene.getObjectByName("Icosphere").rotation.x +=0.1;
    // scene.getObjectByName("Icosphere").rotation.z +=0.01;
    controls.update();
}

function render() {
    renderer.render( scene, camera );
}
animate();