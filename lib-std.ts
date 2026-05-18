// Stand September 2025

// Interne Funktionen

function updMousePos(e:MouseEvent) {
  mouseX = e.offsetX;
  mouseY = e.offsetY;
}

function setMouseDown() {
  mouseStatus = 1;
}

function setMouseUp() {
  mouseStatus = 2;
}

function updTouchPos(e:TouchEvent) {
  e.preventDefault();
  const target = e.target as HTMLElement;
  const touch = e.targetTouches[0];
  const rect = target.getBoundingClientRect();

  mouseX = touch.pageX - rect.left;
  mouseY = touch.pageY - rect.top;

  /* Alte Version
  e.preventDefault();
  //@ts-ignore
  mouseX = e.targetTouches[0].pageX - e.target.getBoundingClientRect().left;
  //@ts-ignore
  mouseY = e.targetTouches[0].pageY - e.target.getBoundingClientRect().top;
  */
}

function setTouchDown(e:TouchEvent) {
  mouseStatus = 1;
  updTouchPos(e);
}

function parseColor(...color:(string|number)[]) {
  let c;
  //Wenn ein Argument übergeben wurde und es ein String ist (Hex-Farbwert)
  if (color.length == 1) {
    if (typeof color[0] == 'string') {
      c = color[0];
    // Wenn ein Argument übergeben wurde und es eine Zahl ist (RGB)
    } else {
      c = `RGB(${color[0]},${color[0]},${color[0]})`
    }
  //Ansonsten 3 RGB Zahlen  
  } else {
    c = `RGB(${color[0]||255},${color[1]||255},${color[2]||255})`
  }
  return c;
}

// Instanzvariablen des Moduls
const canv = document.querySelector("canvas") as HTMLCanvasElement;
const output = document.querySelector("#output") as HTMLLabelElement;
const checkboxes = document.querySelector("#checkbox") as HTMLFieldSetElement;
const ctx = canv.getContext("2d") as CanvasRenderingContext2D;
export let mouseX = 0;
export let mouseY = 0;
let mouseStatus = 0;
let loop = true;

export function init(w:number, h:number) {
  canv.width = w;
  canv.height = h;
  ctx.lineWidth = 1;
  ctx.save();
  canv.addEventListener("mousemove", updMousePos);
  canv.addEventListener("mousedown", setMouseDown);
  canv.addEventListener("mouseup", setMouseUp);
  canv.addEventListener("touchmove", updTouchPos);
  canv.addEventListener("touchstart", setTouchDown);
  canv.addEventListener("touchend", setTouchUp);
}

export function startAnimation(fnDraw:() => void) {
  let draw = fnDraw;
  loop = true;
  let animate = () => {
    draw();
    if (loop) {
      window.requestAnimationFrame(animate);
    }
  };
  window.requestAnimationFrame(animate);
}

export function getWidth():number {
  return canv.width;
}


export function getHeight():number {
  return canv.height;
}

export function noLoop() {
  loop = false;
}

export function setTouchUp() {
  mouseStatus = 2;
}

export function isMouseDown():boolean {
  if (mouseStatus == 1) {
    return true;
  } else {
    return false;
  }
}

export function isMouseUp():boolean {
  if (mouseStatus == 2) {
    mouseStatus = 0;
    return true;
  } else {
    return false;
  }
}

export function createP(item:string) {
  const newItem = document.createTextNode(item);
  output.appendChild(newItem);
}

export function createCheckbox(attribut: string, text: string) {
  const label = document.createElement('label');
  const input = document.createElement('input');

  input.setAttribute('name', attribut);
  input.setAttribute('type', 'checkbox');
  input.setAttribute('role', 'switch'); // Für das Styling, z.B. bei Frameworks

  const labelText = document.createTextNode(text);

  label.appendChild(input);
  label.appendChild(labelText);

  //Das fertige Label zum Container (z.B. zum Body) hinzufügen
  checkboxes.appendChild(label);

}

/** saves the current drawing state. Use together with pop 
*/
export function push() {
  ctx.save();
}

/** restores drawing state. Use together with push 
*/
export function pop() {
  ctx.restore();
}

/** Transformation to current matrix 
 */
export function translate(x:number, y:number) {
  ctx.translate(x, y);
}

/** rotates drawing context in degrees
 */
export function rotate(n:number) {
  ctx.rotate(n);
}

/** Fillcolor of shape
 * color Hex | RGB r, g, b
 */
export function fillColor(...color:(string|number)[]) {
  const c = parseColor(...color);
  ctx.fillStyle = c;
}

/** Stroke Gradiant 
 *  color Hex | RGB r, g, b
 */
export function strokeGrd(x:number, y:number, max:number, ...color:(string|number)[]) {
  const c = parseColor(...color);
  const grd = ctx.createRadialGradient(x, y, 5, x, y, max);
  grd.addColorStop(0, c);
  grd.addColorStop(1, `RGB(${0},${0},${0})`);
  ctx.strokeStyle = grd;
}

/** Strokecolor of shape
 *  Hex | RGB r, g, b
 */
export function strokeColor(...color:(string|number)[]) {
  const c = parseColor(...color);
  ctx.strokeStyle = c;
}

/** strokewidth of line
 */
export function strokeWidth(w:number) {
  ctx.lineWidth = w;
}

/** background color of drawing context
 *  Hex | RGB r, g, b  
 */
export function background(...color:(string|number)[]) {
  const c = parseColor(...color);
  push();
  ctx.fillStyle = c;
  ctx.fillRect(0, 0, canv.width, canv.height);
  pop();
}

/** drawing rectangle with start position
 * position-x 
 * position-y
 * width
 * height
 * style 0, 1 or 2 */
export function rect(x:number, y:number, w:number, h:number, style = 0) {
  if (style == 0) {
    ctx.strokeRect(x, y, w, h);
  }
  if (style == 1) {
    ctx.fillRect(x, y, w, h);
  }
  if (style == 2) {
    ctx.strokeRect(x, y, w, h);
    ctx.fillRect(x, y, w, h);
  }
}

/** drawing line from point1 (x1,y1) to point2 (x2,y2)
 */
export function line(x1:number, y1:number, x2:number, y2:number) {
  const path = new Path2D();
  path.moveTo(x1, y1);
  path.lineTo(x2, y2);
  path.closePath();
  ctx.stroke(path);
}

/** drawing triangle with point1, point2 and point3
 */
export function triangle(x1:number, y1:number, x2:number, y2:number, x3:number, y3:number, style = 0) {
  const path = new Path2D();
  path.moveTo(x1, y1);
  path.lineTo(x2, y2);
  path.lineTo(x3, y3);
  path.closePath();
  if (style == 0) {
    ctx.stroke(path);
  }
  if (style == 1) {
    ctx.fill(path);
  }
  if (style == 2) {
    ctx.stroke(path);
    ctx.fill(path);
  }
}

/** drawing shape with 4 points
 */
export function shape(x1:number, y1:number, x2:number, y2:number, x3:number, y3:number, x4:number, y4:number, style = 0) {
  const path = new Path2D();
  path.moveTo(x1, y1);
  path.lineTo(x2, y2);
  path.lineTo(x3, y3);
  path.lineTo(x4, y4);
  path.closePath();
  if (style == 0) {
    ctx.stroke(path);
  }
  if (style == 1) {
    ctx.fill(path);
  }
  if (style == 2) {
    ctx.stroke(path);
    ctx.fill(path);
  }
}

export function circle(x:number, y:number, radius:number, style = 0) {
  const path = new Path2D();
  path.arc(x, y, radius, 0, 2 * Math.PI);
  path.closePath();
  if (style == 0) {
    ctx.stroke(path);
  }
  if (style == 1) {
    ctx.fill(path);
  }
  if (style == 2) {
    ctx.stroke(path);
    ctx.fill(path);
  }
}

export function drawArrow(v_base:Vector, v_target:Vector, ...myColor:(string|number)[]) {
  const v_heading = subVector(v_target, v_base);
  push();
  strokeColor(...myColor);
  strokeWidth(3);
  fillColor(...myColor);
  translate(v_base.x, v_base.y);
  line(0, 0, v_heading.x, v_heading.y);
  rotate(v_heading.heading());
  let arrowSize = 7;
  translate(v_heading.mag() - arrowSize, 0);
  triangle(0, arrowSize / 2, 0, -arrowSize / 2, arrowSize, 0);
  pop();
}

///////////////////////////////////////////////
//     Ab hier Implementierung für Vector !!!
///////////////////////////////////////////////

export class Vector {
  x: number;
  y: number;

  constructor(x:number, y:number) {
    this.x = x;
    this.y = y;
  }
 
  set(x:number, y:number) {
    this.x = x;
    this.y = y;
  }
    
  copy():Vector {
    return new Vector(this.x, this.y);
  }

  div(n:number) {
    this.x /= n;
    this.y /= n;
  }

  mult(n:number) {
    this.x *= n;
    this.y *= n;
  }

  magSq():number {
    return this.x * this.x + this.y * this.y;
  }

  mag():number {
    return Math.sqrt(this.magSq());
  }

  add(v:Vector) {
    this.x += v.x;
    this.y += v.y;
  }

  sub(v:Vector) {
    this.x -= v.x;
    this.y -= v.y;
  }

  dist(v:Vector):number {
    const vdist = this.copy();
    vdist.sub(v);
    return vdist.mag();
  }

  normalize() {
    const len = this.mag();
    if (len != 0) {
      this.div(len);
    }
  }

  limit(max:number) {
    const mSq = this.magSq();
    if (mSq > max * max) {
      this.setMag(max);
    }
  }
  
  setMag(magnitude:number) {
    this.normalize();
    this.mult(magnitude);
  }
  
  dot(v:Vector):number {
    return this.x * v.x + this.y * v.y;
  }
  
  cross(v:Vector):number {
    return this.x * v.y - this.y * v.x;
  }

  heading():number {
    const h = Math.atan2(this.y, this.x);
    return h;
  }
  
  rotate(base:Vector, n:number) {
    const direction = this.copy();
    direction.sub(base);
    const newHeading = direction.heading() + n;
    const magnitude = direction.mag();
    this.x = base.x + Math.cos(newHeading) * magnitude;
    this.y = base.y + Math.sin(newHeading) * magnitude;
  }
  
  rotateMatrix(base:Vector, n:number) {
    const direction = this.copy();
    direction.sub(base);
    const x = direction.x * Math.cos(n) - direction.y * Math.sin(n);
    const y = direction.x * Math.sin(n) + direction.y * Math.cos(n);
    this.x = x + base.x;
    this.y = y + base.y;
  }

  angleBetween(v:Vector):number {
    const dotmagmag = this.dot(v) / (this.mag() * v.mag());
    const angle = Math.acos(Math.min(1, Math.max(-1, dotmagmag)));
    return angle;
  }

  perp():Vector {
    return new Vector(-this.y, this.x);
  }
}

export function VectorRandom2D():Vector {
  const angle = Math.random() * Math.PI * 2;
  const x = Math.cos(angle);
  const y = Math.sin(angle);
  return new Vector(x, y);
}

export function fromAngle(angle:number, len:number):Vector {
  const x = len * Math.cos(angle);
  const y = len * Math.sin(angle);
  return new Vector(x, y);
}

export function addVector(v1:Vector, v2:Vector) {
  return new Vector(v1.x + v2.x, v1.y + v2.y);
}

export function subVector(v1:Vector, v2:Vector) {
  return new Vector(v1.x - v2.x, v1.y - v2.y);
}

export function multVector(v:Vector, n:number):Vector {
  const vmult = v.copy();
  vmult.mult(n);
  return vmult;
}

export function divVector(v:Vector, n:number):Vector {
  const vdiv = v.copy();
  vdiv.div(n);
  return vdiv;
}

export function normVector(a: Vector): Vector {
  const v = a.copy();
  v.normalize();
  return v;
}

export function dotProduct(v1:Vector, v2:Vector):number {
  return v1.dot(v2);
}

export function crossProduct(v1:Vector, v2:Vector):number {
  return v1.cross(v2);
}

/** Returns Point of intersection + Scalar s between 
 * line a (a0->a1) and line b (b0->b1)
 * Example: const [Point, s] = intersect(a0, a1, b0, b1)
  */
 export function intersect(a0:Vector, a1:Vector, b0:Vector, b1:Vector):[Vector|null, number|null] {
  let pt:[Vector|null, number|null]  = [null, null]; // pt[0] = Point of intersection, pt[1] = s
  const a = subVector(a1, a0);
  const b = subVector(b1, b0);
  const den1 = a.cross(b);
  const den2 = b.cross(a);

  if (den1 != 0) {
    const s = subVector(b0, a0).cross(b) / den1;
    const u = subVector(a0, b0).cross(a) / den2;
    if (s > 0 && s < 1 && u > 0 && u < 1) {
      const p = addVector(a0, multVector(a, s));
      pt[0] = p;
      pt[1] = s;
    }
  }

  return pt;
}

/** Mindistance between point p and line a(a0->a1)
 */
export function minDist(p:Vector, a0:Vector, a1:Vector):[number, Vector] {
  let dist = -1;
  let normalPoint = new Vector(0, 0);

  //Vektor line a0 to a1
  const a0a1 = subVector(a1, a0);
  //Vektor line a0 to p
  const a0p = subVector(p, a0);
  //Magnitude of line a0 to a1
  const magnitude = a0a1.mag();

  //Scalarproduct from line a0p to line a0a1
  a0a1.normalize();
  const sp = a0p.dot(a0a1);

  //Scalarprojection in magnitude of line a0a1?
  if (sp > 0 && sp <= magnitude) {
    a0a1.mult(sp);
    normalPoint = addVector(a0, a0a1);
    dist = p.dist(normalPoint);
  }
  return [dist, normalPoint];
}