// lib-phys.ts – Port von phys.go (Stand April 2026)
// Verwendet std.Vector aus lib-std.ts statt raylib Vector2

import * as std from "./lib-2d.ts";

const INF = Number.MAX_VALUE;

// ─── Interfaces ──────────────────────────────────────────────────────────────

export interface Shape {
  update(dt: number): void;
  draw(color: string, thick: number): void;
  rotate(angle: number): void;
  applyForce(force: std.Vector, angForce: number): void;
  resetPos(delta: std.Vector): void;
  applyGravity(gravity: std.Vector): void;
  clearState(): void;
}

interface BasicShape {
  location: std.Vector;
  velocity: std.Vector;
  angVelocity: number;
  accel: std.Vector;
  angAccel: number;
  mass: number;
  inertia: number;
  isGrounded: boolean;
}

// ─── Polygon ─────────────────────────────────────────────────────────────────

export class Polygon implements Shape {
  basic: BasicShape;
  vertices: std.Vector[];

  constructor(x: number, y: number, w: number, h: number, wall: boolean) {
    const mass = wall ? INF : w * h;
    const inertia = wall ? INF : (mass * (w * w + h * h)) / 2;

    this.basic = {
      location: new std.Vector(x + w / 2, y + h / 2),
      velocity: new std.Vector(0, 0),
      angVelocity: 0,
      accel: new std.Vector(0, 0),
      angAccel: 0,
      mass,
      inertia,
      isGrounded: false,
    };

    this.vertices = [
      new std.Vector(x, y),
      new std.Vector(x + w, y),
      new std.Vector(x + w, y + h),
      new std.Vector(x, y + h),
    ];
  }

  applyForce(force: std.Vector, angForce: number) {
    this.basic.accel = std.addVector(this.basic.accel, force);
    this.basic.angAccel += angForce;
  }

  update(dt: number) {
    this.basic.velocity = std.addVector(
      this.basic.velocity,
      std.multVector(this.basic.accel, dt),
    );
    this.basic.angVelocity += this.basic.angAccel * dt;

    this.basic.velocity = std.multVector(this.basic.velocity, 0.9995);
    this.basic.location = std.addVector(
      this.basic.location,
      std.multVector(this.basic.velocity, dt),
    );
    for (let i = 0; i < this.vertices.length; i++) {
      this.vertices[i] = std.addVector(
        this.vertices[i],
        std.multVector(this.basic.velocity, dt),
      );
    }
    this.basic.angVelocity *= 0.9995;
    this.rotate(this.basic.angVelocity * dt);

    this.basic.accel = new std.Vector(0, 0);
    this.basic.angAccel = 0;
  }

  draw(color: string, thick: number) {
    const n = this.vertices.length;
    std.strokeColor(color);
    std.strokeWidth(thick);
    for (let i = 0; i < n; i++) {
      const a = this.vertices[i];
      const b = this.vertices[(i + 1) % n];
      std.line(a.x, a.y, b.x, b.y);
    }
    // Mittelpunkt
    std.fillColor(color);
    std.circle(this.basic.location.x, this.basic.location.y, 3, 1);
  }

  rotate(angle: number) {
    for (let i = 0; i < this.vertices.length; i++) {
      this.vertices[i].rotateMatrix(this.basic.location, angle);
    }
  }

  resetPos(delta: std.Vector) {
    for (let i = 0; i < this.vertices.length; i++) {
      this.vertices[i] = std.addVector(this.vertices[i], delta);
    }
    this.basic.location = std.addVector(this.basic.location, delta);
  }

  applyGravity(gravity: std.Vector) {
    if (this.basic.mass < INF) {
      if (!this.basic.isGrounded) {
        this.applyForce(gravity, 0);
      } else {
        if (
          this.basic.velocity.mag() < 0.5 &&
          Math.abs(this.basic.angVelocity) < 0.1
        ) {
          this.basic.velocity = new std.Vector(0, 0);
          this.basic.angVelocity = 0;
        } else {
          const dampingForce = std.multVector(this.basic.velocity, -0.5);
          const dampingAngForce = this.basic.angVelocity * -0.5;
          this.applyForce(dampingForce, dampingAngForce);
        }
      }
    }
  }

  clearState() {
    this.basic.isGrounded = false;
  }
}

// ─── Circle ──────────────────────────────────────────────────────────────────

export class Circle implements Shape {
  basic: BasicShape;
  radius: number;
  orientation: std.Vector; // Punkt auf Kreisrand (für Rotationsanzeige)

  constructor(x: number, y: number, r: number, wall: boolean) {
    const mass = wall ? INF : r * r * 2;
    const inertia = wall ? INF : r * r * r * 100;

    this.basic = {
      location: new std.Vector(x, y),
      velocity: new std.Vector(0, 0),
      angVelocity: 0,
      accel: new std.Vector(0, 0),
      angAccel: 0,
      mass,
      inertia,
      isGrounded: false,
    };
    this.radius = r;
    this.orientation = new std.Vector(r + x, y);
  }

  applyForce(force: std.Vector, angForce: number) {
    this.basic.accel = std.addVector(this.basic.accel, force);
    this.basic.angAccel += angForce;
  }

  update(dt: number) {
    this.basic.velocity = std.addVector(
      this.basic.velocity,
      std.multVector(this.basic.accel, dt),
    );
    this.basic.angVelocity += this.basic.angAccel * dt;

    this.basic.velocity = std.multVector(this.basic.velocity, 0.9995);
    this.basic.location = std.addVector(
      this.basic.location,
      std.multVector(this.basic.velocity, dt),
    );
    this.orientation = std.addVector(
      this.orientation,
      std.multVector(this.basic.velocity, dt),
    );

    this.basic.angVelocity *= 0.9995;
    this.rotate(this.basic.angVelocity * dt);

    this.basic.accel = new std.Vector(0, 0);
    this.basic.angAccel = 0;
  }

  draw(color: string, thick: number) {
    std.strokeColor(color);
    std.strokeWidth(thick);
    std.circle(this.basic.location.x, this.basic.location.y, this.radius, 0);
    // Orientierungslinie
    std.line(
      this.basic.location.x,
      this.basic.location.y,
      this.orientation.x,
      this.orientation.y,
    );
    // Mittelpunkt
    std.fillColor(color);
    std.circle(this.basic.location.x, this.basic.location.y, 3, 1);
  }

  rotate(angle: number) {
    this.orientation.rotateMatrix(this.basic.location, angle);
  }

  resetPos(delta: std.Vector) {
    this.basic.location = std.addVector(this.basic.location, delta);
    this.orientation = std.addVector(this.orientation, delta);
  }

  applyGravity(gravity: std.Vector) {
    if (this.basic.mass < INF) {
      if (!this.basic.isGrounded) {
        this.applyForce(gravity, 0);
      } else {
        const dampingForce = std.multVector(this.basic.velocity, -0.5);
        const dampingAngForce = this.basic.angVelocity * -0.5;
        this.applyForce(dampingForce, dampingAngForce);
      }
    }
  }

  clearState() {
    this.basic.isGrounded = false;
  }
}

// ─── Interne Geometrie-Hilfsfunktionen ───────────────────────────────────────

function projectPolygon(
  vertices: std.Vector[],
  axis: std.Vector,
): [number, number] {
  if (vertices.length === 0) return [0, 0];
  let min = std.dotProduct(vertices[0], axis);
  let max = min;
  for (let i = 1; i < vertices.length; i++) {
    const proj = std.dotProduct(vertices[i], axis);
    if (proj < min) min = proj;
    if (proj > max) max = proj;
  }
  return [min, max];
}

function getAxes(vertices: std.Vector[]): std.Vector[] {
  const axes: std.Vector[] = [];
  const n = vertices.length;
  for (let i = 0; i < n; i++) {
    const a = vertices[i];
    const b = vertices[(i + 1) % n];
    const edge = std.normVector(std.subVector(b, a));
    axes.push(new std.Vector(-edge.y, edge.x));
  }
  return axes;
}

function pointInPolygon(point: std.Vector, vertices: std.Vector[]): boolean {
  const n = vertices.length;
  for (let i = 0; i < n; i++) {
    const a = vertices[i];
    const b = vertices[(i + 1) % n];
    const edge = std.subVector(b, a);
    const toPoint = std.subVector(point, a);
    if (std.dotProduct(edge, toPoint) < 0) return false;
  }
  return true;
}

function findContactPoints(
  verticesA: std.Vector[],
  verticesB: std.Vector[],
): std.Vector[] {
  const contacts: std.Vector[] = [];
  for (const p of verticesA) {
    if (pointInPolygon(p, verticesB)) contacts.push(p);
  }
  for (const p of verticesB) {
    if (pointInPolygon(p, verticesA)) contacts.push(p);
  }
  return contacts;
}

function findReferenceEdge(
  vertices: std.Vector[],
  normal: std.Vector,
): [std.Vector, std.Vector] {
  let bestDot = -INF;
  let p1 = new std.Vector(0, 0);
  let p2 = new std.Vector(0, 0);
  const n = vertices.length;
  for (let i = 0; i < n; i++) {
    const a = vertices[i];
    const b = vertices[(i + 1) % n];
    const edge = std.normVector(std.subVector(b, a));
    const edgeNormal = new std.Vector(-edge.y, edge.x);
    const d = std.dotProduct(edgeNormal, normal);
    if (d > bestDot) {
      bestDot = d;
      p1 = a;
      p2 = b;
    }
  }
  return [p1, p2];
}

function projectPointOntoEdge(
  p: std.Vector,
  a: std.Vector,
  b: std.Vector,
): std.Vector {
  const ab = std.subVector(b, a);
  let t = std.dotProduct(std.subVector(p, a), ab) / std.dotProduct(ab, ab);
  t = Math.min(1, Math.max(0, t));
  return std.addVector(a, std.multVector(ab, t));
}

function transferContactsToA(
  contacts: std.Vector[],
  verticesA: std.Vector[],
  normal: std.Vector,
): std.Vector[] {
  const [refStart, refEnd] = findReferenceEdge(verticesA, normal);
  return contacts.map((c) => projectPointOntoEdge(c, refStart, refEnd));
}

// ─── Massen-basierte Positions-Resets ────────────────────────────────────────

function resetPolyPositionsBasedOnMass(
  polyA: Polygon,
  polyB: Polygon,
  mtv: std.Vector,
) {
  if (polyA.basic.mass < INF) {
    polyA.resetPos(std.multVector(mtv, -0.5));
  } else {
    polyB.resetPos(std.multVector(mtv, 0.5));
  }
  if (polyB.basic.mass < INF) {
    polyB.resetPos(std.multVector(mtv, 0.5));
  } else {
    polyA.resetPos(std.multVector(mtv, -0.5));
  }
}

function resetCirclePolyPositionsBasedOnMass(
  poly: Polygon,
  circle: Circle,
  mtv: std.Vector,
) {
  if (poly.basic.mass < INF) {
    poly.resetPos(std.multVector(mtv, -0.5));
  } else {
    circle.resetPos(std.multVector(mtv, 0.5));
  }
  if (circle.basic.mass < INF) {
    circle.resetPos(std.multVector(mtv, 0.5));
  } else {
    poly.resetPos(std.multVector(mtv, -0.5));
  }
}

function resetCirclePositionsBasedOnMass(
  cA: Circle,
  cB: Circle,
  mtv: std.Vector,
) {
  if (cA.basic.mass < INF) {
    cA.resetPos(std.multVector(mtv, -0.5));
  } else {
    cB.resetPos(std.multVector(mtv, 0.5));
  }
  if (cB.basic.mass < INF) {
    cB.resetPos(std.multVector(mtv, 0.5));
  } else {
    cA.resetPos(std.multVector(mtv, -0.5));
  }
}

// ─── Relative Vektoren (für Impulsberechnung) ────────────────────────────────

type RelVectors = {
  rAP_perp: std.Vector;
  rBP_perp: std.Vector;
  vGesamtA: std.Vector;
  velocity_AB: std.Vector;
};

function calcPolyRelVectors(
  polyA: Polygon,
  polyB: Polygon,
  cp: std.Vector,
): RelVectors {
  const rAP = std.subVector(cp, polyA.basic.location);
  const rBP = std.subVector(cp, polyB.basic.location);
  const rAP_perp = new std.Vector(-rAP.y, rAP.x);
  const rBP_perp = new std.Vector(-rBP.y, rBP.x);
  const VtanA = std.multVector(rAP_perp, polyA.basic.angVelocity);
  const VtanB = std.multVector(rBP_perp, polyB.basic.angVelocity);
  const vGesamtA = std.addVector(polyA.basic.velocity, VtanA);
  const vGesamtB = std.addVector(polyB.basic.velocity, VtanB);
  return {
    rAP_perp,
    rBP_perp,
    vGesamtA,
    velocity_AB: std.subVector(vGesamtA, vGesamtB),
  };
}

function calcCirclePolyRelVectors(
  poly: Polygon,
  circle: Circle,
  cp: std.Vector,
): RelVectors {
  const rAP = std.subVector(cp, poly.basic.location);
  const rBP = std.subVector(cp, circle.basic.location);
  const rAP_perp = new std.Vector(-rAP.y, rAP.x);
  const rBP_perp = new std.Vector(-rBP.y, rBP.x);
  const VtanA = std.multVector(rAP_perp, poly.basic.angVelocity);
  const VtanB = std.multVector(rBP_perp, circle.basic.angVelocity);
  const vGesamtA = std.addVector(poly.basic.velocity, VtanA);
  const vGesamtB = std.addVector(circle.basic.velocity, VtanB);
  return {
    rAP_perp,
    rBP_perp,
    vGesamtA,
    velocity_AB: std.subVector(vGesamtA, vGesamtB),
  };
}

function calcCircleRelVectors(
  cA: Circle,
  cB: Circle,
  mtv: std.Vector,
): RelVectors {
  const nMtv = std.normVector(mtv);
  const rA = std.multVector(nMtv, -cA.radius);
  const rB = std.multVector(nMtv, cB.radius);
  const rAP_perp = new std.Vector(-rA.y, rA.x);
  const rBP_perp = new std.Vector(-rB.y, rB.x);
  const VtanA = std.multVector(rAP_perp, cA.basic.angVelocity);
  const VtanB = std.multVector(rBP_perp, cB.basic.angVelocity);
  const vGesamtA = std.addVector(cA.basic.velocity, VtanA);
  const vGesamtB = std.addVector(cB.basic.velocity, VtanB);
  return {
    rAP_perp,
    rBP_perp,
    vGesamtA,
    velocity_AB: std.subVector(vGesamtA, vGesamtB),
  };
}

// ─── Grounded-Checks ─────────────────────────────────────────────────────────

function isHorizontal(mtv: std.Vector, gravity: std.Vector): boolean {
  const result = (std.dotProduct(mtv, gravity) / mtv.mag()) * gravity.mag();
  return result > 0.9 || result < -0.9;
}

function isCircleGrounded(
  angVel: number,
  velocity_AB: std.Vector,
  mtv: std.Vector,
): boolean {
  if (velocity_AB.mag() > 0.5) return false;
  if (Math.abs(angVel) > 0.1) return false;
  const gravity = new std.Vector(0, 1);
  const dot = std.dotProduct(mtv, gravity);
  const lenM = mtv.mag();
  const lenG = gravity.mag();
  if (lenM === 0 || lenG === 0) return false;
  return dot / (lenM * lenG) < -0.9;
}

function isPolyGrounded(
  contacts: std.Vector[],
  vGesamtA: std.Vector,
  velocity_AB: std.Vector,
  mtv: std.Vector,
): boolean {
  return (
    contacts.length > 1 &&
    velocity_AB.mag() < 1.5 &&
    vGesamtA.mag() < 1.5 &&
    isHorizontal(mtv, new std.Vector(0, 1))
  );
}

function findPolyTop(polyA: Polygon, polyB: Polygon, mtv: std.Vector): Polygon {
  return std.dotProduct(mtv, new std.Vector(0, 1)) > 0 ? polyA : polyB;
}

function isCenterOfMassSupported(
  location: std.Vector,
  contacts: std.Vector[],
  normal: std.Vector,
): boolean {
  if (contacts.length < 2) return false;
  const tangent = new std.Vector(-normal.y, normal.x);
  const comProj = std.dotProduct(location, tangent);
  let minP = INF;
  let maxP = -INF;
  for (const cp of contacts) {
    const proj = std.dotProduct(cp, tangent);
    if (proj < minP) minP = proj;
    if (proj > maxP) maxP = proj;
  }
  const epsilon = 0.5;
  return comProj >= minP - epsilon && comProj <= maxP + epsilon;
}

// ─── Impuls- und Kraftberechnungen ───────────────────────────────────────────

function calcFrictionVector(
  mtv: std.Vector,
  velocity_AB: std.Vector,
): std.Vector {
  const t = new std.Vector(-mtv.y, mtv.x);
  const sp = std.dotProduct(velocity_AB, t);
  return std.normVector(std.multVector(t, sp));
}

function calcPolyImpulse(
  polyA: Polygon,
  polyB: Polygon,
  mtv: std.Vector,
  rAP_perp: std.Vector,
  rBP_perp: std.Vector,
  velocity_AB: std.Vector,
  e: number,
): number {
  const jNum = std.dotProduct(std.multVector(velocity_AB, -(1 + e)), mtv);
  const jLinear = std.dotProduct(
    mtv,
    std.multVector(mtv, 1 / polyA.basic.mass + 1 / polyB.basic.mass),
  );
  const jAng =
    Math.pow(std.dotProduct(rAP_perp, mtv), 2) / polyA.basic.inertia +
    Math.pow(std.dotProduct(rBP_perp, mtv), 2) / polyB.basic.inertia;
  return jNum / (jLinear + jAng);
}

function calcCirclePolyImpulse(
  poly: Polygon,
  circle: Circle,
  mtv: std.Vector,
  rAP_perp: std.Vector,
  velocity_AB: std.Vector,
  e: number,
): number {
  const jNum = std.dotProduct(std.multVector(velocity_AB, -(1 + e)), mtv);
  const jLinear = std.dotProduct(
    mtv,
    std.multVector(mtv, 1 / poly.basic.mass + 1 / circle.basic.mass),
  );
  const jAng = Math.pow(std.dotProduct(rAP_perp, mtv), 2) / poly.basic.inertia;
  return jNum / (jLinear + jAng);
}

function calcCircleImpulse(
  cA: Circle,
  cB: Circle,
  mtv: std.Vector,
  velocity_AB: std.Vector,
  e: number,
): number {
  const jNum = std.dotProduct(std.multVector(velocity_AB, -(1 + e)), mtv);
  const jLinear = std.dotProduct(
    mtv,
    std.multVector(mtv, 1 / cA.basic.mass + 1 / cB.basic.mass),
  );
  return jNum / jLinear;
}

type Forces = {
  forceA: std.Vector;
  angForceA: number;
  forceB: std.Vector;
  angForceB: number;
};

function calcPolyCollisionForces(
  polyA: Polygon,
  polyB: Polygon,
  mtv: std.Vector,
  rAP_perp: std.Vector,
  rBP_perp: std.Vector,
  velocity_AB: std.Vector,
): Forces {
  const e = 0.4;
  const t = calcFrictionVector(mtv, velocity_AB);
  const f = -0.15;
  const j = calcPolyImpulse(
    polyA,
    polyB,
    mtv,
    rAP_perp,
    rBP_perp,
    velocity_AB,
    e,
  );

  const forceA = std.addVector(
    std.multVector(mtv, j / polyA.basic.mass),
    std.multVector(t, (f * -j) / polyA.basic.mass),
  );
  const angForceA = std.dotProduct(
    rAP_perp,
    std.addVector(
      std.multVector(mtv, j / polyA.basic.inertia),
      std.multVector(t, (f * -j) / polyA.basic.inertia),
    ),
  );
  const forceB = std.addVector(
    std.multVector(mtv, -j / polyB.basic.mass),
    std.multVector(t, (f * j) / polyB.basic.mass),
  );
  const angForceB = std.dotProduct(
    rBP_perp,
    std.addVector(
      std.multVector(mtv, -j / polyB.basic.inertia),
      std.multVector(t, (f * j) / polyB.basic.inertia),
    ),
  );
  return { forceA, angForceA, forceB, angForceB };
}

function calcCirclePolyCollisionForces(
  poly: Polygon,
  circle: Circle,
  mtv: std.Vector,
  rAP_perp: std.Vector,
  rBP_perp: std.Vector,
  velocity_AB: std.Vector,
): Forces {
  const e = 0.3;
  const t = calcFrictionVector(mtv, velocity_AB);
  const f = -0.15;
  const j = calcCirclePolyImpulse(poly, circle, mtv, rAP_perp, velocity_AB, e);

  const forceA = std.addVector(
    std.multVector(mtv, j / poly.basic.mass),
    std.multVector(t, (f * -j) / poly.basic.mass),
  );
  const angForceA = std.dotProduct(
    rAP_perp,
    std.addVector(
      std.multVector(mtv, j / poly.basic.inertia),
      std.multVector(t, (f * -j) / poly.basic.inertia),
    ),
  );
  const forceB = std.addVector(
    std.multVector(mtv, -j / circle.basic.mass),
    std.multVector(t, (f * j) / circle.basic.mass),
  );
  const angForceB = std.dotProduct(
    rBP_perp,
    std.multVector(t, (f * j) / circle.basic.inertia),
  );
  return { forceA, angForceA, forceB, angForceB };
}

function calcCircleCollisionForces(
  cA: Circle,
  cB: Circle,
  mtv: std.Vector,
  rAP_perp: std.Vector,
  rBP_perp: std.Vector,
  velocity_AB: std.Vector,
): Forces {
  const e = 0.3;
  const t = calcFrictionVector(mtv, velocity_AB);
  const f = -0.15;
  const j = calcCircleImpulse(cA, cB, mtv, velocity_AB, e);

  const forceA = std.addVector(
    std.multVector(mtv, j / cA.basic.mass),
    std.multVector(t, (f * j) / cA.basic.mass),
  );
  const angForceA = std.dotProduct(
    rAP_perp,
    std.multVector(t, (f * j) / cA.basic.inertia),
  );
  const forceB = std.addVector(
    std.multVector(mtv, -j / cB.basic.mass),
    std.multVector(t, (f * -j) / cB.basic.mass),
  );
  const angForceB = std.dotProduct(
    rBP_perp,
    std.multVector(t, (f * -j) / cB.basic.inertia),
  );
  return { forceA, angForceA, forceB, angForceB };
}

// ─── Kollisionserkennung ─────────────────────────────────────────────────────

type CollResult = {
  isColliding: boolean;
  mtv: std.Vector;
  contacts: std.Vector[];
};
const NO_COLL: CollResult = {
  isColliding: false,
  mtv: new std.Vector(0, 0),
  contacts: [],
};

export function detectCollisionPoly(
  polyA: Polygon,
  polyB: Polygon,
): CollResult {
  let smallestOverlap = INF;
  let smallestAxis = new std.Vector(0, 0);

  const axes = [...getAxes(polyA.vertices), ...getAxes(polyB.vertices)];
  for (const axis of axes) {
    const [minA, maxA] = projectPolygon(polyA.vertices, axis);
    const [minB, maxB] = projectPolygon(polyB.vertices, axis);
    const overlap = Math.min(maxA, maxB) - Math.max(minA, minB);
    if (overlap <= 0) return NO_COLL;

    if (overlap < smallestOverlap) {
      smallestOverlap = overlap;
      smallestAxis = axis;
      const dir = std.subVector(polyB.basic.location, polyA.basic.location);
      if (std.dotProduct(dir, smallestAxis) < 0)
        smallestAxis = std.multVector(smallestAxis, -1);
    }
  }

  const mtv = std.multVector(smallestAxis, smallestOverlap);
  const rawContacts = findContactPoints(polyA.vertices, polyB.vertices);
  const contacts = transferContactsToA(
    rawContacts,
    polyA.vertices,
    std.multVector(smallestAxis, -1),
  );
  if (contacts.length > 0) return { isColliding: true, mtv, contacts };
  return NO_COLL;
}

export function detectCollisionCirclePoly(
  poly: Polygon,
  circle: Circle,
): CollResult {
  const n = poly.vertices.length;
  let closestDistSq = INF;
  let closestPoint = new std.Vector(0, 0);
  let inside = true;

  for (let i = 0; i < n; i++) {
    const a = poly.vertices[i];
    const b = poly.vertices[(i + 1) % n];
    const edge = std.subVector(b, a);
    const toCircle = std.subVector(circle.basic.location, a);
    const normal = new std.Vector(-edge.y, edge.x);

    if (std.dotProduct(normal, toCircle) < 0) inside = false;

    const edgeLenSq = std.dotProduct(edge, edge);
    const t = std.dotProduct(toCircle, edge) / edgeLenSq;
    const clamped = Math.min(1, Math.max(0, t));
    const current = std.addVector(a, std.multVector(edge, clamped));
    const diff = std.subVector(circle.basic.location, current);
    const distSq = std.dotProduct(diff, diff);

    if (distSq < closestDistSq) {
      closestDistSq = distSq;
      closestPoint = current;
    }
  }

  if (inside || closestDistSq <= circle.radius * circle.radius) {
    let axis: std.Vector;
    let overlap: number;

    if (inside) {
      let minOverlap = INF;
      axis = new std.Vector(0, -1);
      for (let i = 0; i < n; i++) {
        const a = poly.vertices[i];
        const b = poly.vertices[(i + 1) % n];
        const edge = std.subVector(b, a);
        const normal = std.normVector(new std.Vector(-edge.y, edge.x));
        const toCircle = std.subVector(circle.basic.location, a);
        const dist = std.dotProduct(normal, toCircle);
        const ov = dist + circle.radius;
        if (ov < minOverlap) {
          minOverlap = ov;
          axis = std.multVector(normal, -1);
        }
      }
      overlap = minOverlap;
    } else {
      const dist = Math.sqrt(closestDistSq);
      overlap = circle.radius - dist;
      if (dist > 1e-6) {
        axis = std.multVector(
          std.subVector(circle.basic.location, closestPoint),
          1 / dist,
        );
      } else {
        axis = new std.Vector(0, -1);
      }
    }

    const mtv = std.multVector(axis, overlap);
    return { isColliding: true, mtv, contacts: [closestPoint] };
  }

  return NO_COLL;
}

export function detectCollisionCircle(cA: Circle, cB: Circle): CollResult {
  const line = std.subVector(cA.basic.location, cB.basic.location);
  const dist = line.mag();
  const radiusSum = cA.radius + cB.radius;
  if (dist < radiusSum) {
    const overlap = dist - radiusSum;
    const mtv = std.multVector(std.normVector(line), overlap);
    return { isColliding: true, mtv, contacts: [] };
  }
  return NO_COLL;
}

// ─── Kollisionsauflösung ──────────────────────────────────────────────────────

export function resolveCollisionPoly(
  polyA: Polygon,
  polyB: Polygon,
  contacts: std.Vector[],
  mtv: std.Vector,
): Forces {
  const slop = 0.5;
  const mtvLength = mtv.mag();

  if (mtvLength > slop) {
    const mtvCorr = std.multVector(std.normVector(mtv), mtvLength - slop);
    resetPolyPositionsBasedOnMass(polyA, polyB, mtvCorr);
  }

  const mtvN = std.normVector(mtv);

  const centerCP =
    contacts.length > 1
      ? std.multVector(std.addVector(contacts[0], contacts[1]), 0.5)
      : contacts[0];

  const { vGesamtA: vGA, velocity_AB: vAB } = calcPolyRelVectors(
    polyA,
    polyB,
    centerCP,
  );

  const polyTop = findPolyTop(polyA, polyB, mtvN);
  if (isPolyGrounded(contacts, vGA, vAB, mtvN)) {
    polyTop.basic.isGrounded = isCenterOfMassSupported(
      polyTop.basic.location,
      contacts,
      mtvN,
    );
  } else {
    polyTop.basic.isGrounded = false;
  }

  let sumForceA = new std.Vector(0, 0);
  let sumForceB = new std.Vector(0, 0);
  let sumAngForceA = 0;
  let sumAngForceB = 0;
  let applied = 0;

  for (const cp of contacts) {
    const { rAP_perp, rBP_perp, velocity_AB } = calcPolyRelVectors(
      polyA,
      polyB,
      cp,
    );
    if (std.dotProduct(velocity_AB, std.multVector(mtvN, -1)) < 0) {
      const f = calcPolyCollisionForces(
        polyA,
        polyB,
        mtvN,
        rAP_perp,
        rBP_perp,
        velocity_AB,
      );
      sumForceA = std.addVector(sumForceA, f.forceA);
      sumAngForceA += f.angForceA;
      sumForceB = std.addVector(sumForceB, f.forceB);
      sumAngForceB += f.angForceB;
      applied++;
    }
  }

  if (applied > 0) {
    return {
      forceA: std.multVector(sumForceA, 1 / applied),
      angForceA: sumAngForceA / applied,
      forceB: std.multVector(sumForceB, 1 / applied),
      angForceB: sumAngForceB / applied,
    };
  }
  return {
    forceA: new std.Vector(0, 0),
    angForceA: 0,
    forceB: new std.Vector(0, 0),
    angForceB: 0,
  };
}

export function resolveCollisionCirclePoly(
  poly: Polygon,
  circle: Circle,
  contacts: std.Vector[],
  mtv: std.Vector,
): Forces {
  const slop = 0.5;
  const mtvLength = mtv.mag();

  if (mtvLength > slop) {
    const mtvCorr = std.multVector(std.normVector(mtv), mtvLength - slop);
    resetCirclePolyPositionsBasedOnMass(poly, circle, mtvCorr);
  }

  const mtvN = std.normVector(mtv);
  const cp = contacts[0];
  const { rAP_perp, rBP_perp, velocity_AB } = calcCirclePolyRelVectors(
    poly,
    circle,
    cp,
  );

  circle.basic.isGrounded = isCircleGrounded(
    circle.basic.angVelocity,
    std.multVector(velocity_AB, -1),
    std.multVector(mtvN, -1),
  );

  if (std.dotProduct(velocity_AB, std.multVector(mtvN, -1)) < 0) {
    return calcCirclePolyCollisionForces(
      poly,
      circle,
      mtvN,
      rAP_perp,
      rBP_perp,
      velocity_AB,
    );
  }
  return {
    forceA: new std.Vector(0, 0),
    angForceA: 0,
    forceB: new std.Vector(0, 0),
    angForceB: 0,
  };
}

export function resolveCollisionCircle(
  cA: Circle,
  cB: Circle,
  mtv: std.Vector,
): Forces {
  const slop = 0.5;
  const mtvLength = mtv.mag();

  if (mtvLength > slop) {
    const mtvCorr = std.multVector(std.normVector(mtv), mtvLength - slop);
    resetCirclePositionsBasedOnMass(cA, cB, mtvCorr);
  }

  const mtvN = std.normVector(mtv);
  const { rAP_perp, rBP_perp, velocity_AB } = calcCircleRelVectors(
    cA,
    cB,
    mtvN,
  );

  cA.basic.isGrounded = isCircleGrounded(
    cA.basic.angVelocity,
    velocity_AB,
    mtvN,
  );
  cB.basic.isGrounded = isCircleGrounded(
    cB.basic.angVelocity,
    std.multVector(velocity_AB, -1),
    std.multVector(mtvN, -1),
  );

  if (std.dotProduct(velocity_AB, std.multVector(mtvN, -1)) < 0) {
    return calcCircleCollisionForces(
      cA,
      cB,
      mtvN,
      rAP_perp,
      rBP_perp,
      velocity_AB,
    );
  }
  return {
    forceA: new std.Vector(0, 0),
    angForceA: 0,
    forceB: new std.Vector(0, 0),
    angForceB: 0,
  };
}

// ─── Öffentliche Dispatcher ───────────────────────────────────────────────────

export function detectCollision(a: Shape, b: Shape): CollResult {
  if (a instanceof Polygon && b instanceof Polygon)
    return detectCollisionPoly(a, b);
  if (a instanceof Polygon && b instanceof Circle)
    return detectCollisionCirclePoly(a, b);
  if (a instanceof Circle && b instanceof Polygon)
    return detectCollisionCirclePoly(b, a);
  if (a instanceof Circle && b instanceof Circle)
    return detectCollisionCircle(a, b);
  return NO_COLL;
}

export function resolveCollision(
  a: Shape,
  b: Shape,
  contacts: std.Vector[],
  mtv: std.Vector,
): Forces {
  if (a instanceof Polygon && b instanceof Polygon)
    return resolveCollisionPoly(a, b, contacts, mtv);
  if (a instanceof Polygon && b instanceof Circle)
    return resolveCollisionCirclePoly(a, b, contacts, mtv);
  if (a instanceof Circle && b instanceof Polygon)
    return resolveCollisionCirclePoly(b, a, contacts, mtv);
  if (a instanceof Circle && b instanceof Circle)
    return resolveCollisionCircle(a, b, mtv);
  return {
    forceA: new std.Vector(0, 0),
    angForceA: 0,
    forceB: new std.Vector(0, 0),
    angForceB: 0,
  };
}
